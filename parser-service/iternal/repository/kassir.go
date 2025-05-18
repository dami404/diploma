package repository

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/dami404/diploma-parser/iternal/entity"
	"github.com/playwright-community/playwright-go"
)

// parsePriceKassir parses the price string into an integer.
func parsePriceKassir(priceStr string) (int, error) {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(priceStr, "\u00a0", ""), " ", "")
	return strconv.Atoi(cleaned)
}

func (r *HTTPDBRepository) ParseKassir(ctx context.Context, city string, name string) []entity.Event {
	select {
	case <-ctx.Done():
		log.Println("Repository.ProfitEvents.ParseKassir: timeout")
		return nil
	default:
		cityToUrl := map[string]string{
			"msk": "https://msk.kassir.ru",
			"spb": "https://spb.kassir.ru",
			"ekb": "https://ekb.kassir.ru",
		}
		if _, ok := cityToUrl[city]; !ok {
			log.Printf("City %s not supported", city)
			return nil
		}

		query := cityToUrl[city] + "/category?keyword=" + name
		log.Println("Parsing query:", query)

		pw, err := playwright.Run()
		if err != nil {
			log.Fatalf("Could not start Playwright: %v", err)
			return nil
		}
		defer pw.Stop()

		browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true),
		})
		if err != nil {
			log.Printf("Could not launch browser: %v", err)
			return nil
		}
		defer browser.Close()

		page, err := browser.NewPage()
		if err != nil {
			log.Printf("Could not create page: %v", err)
			return nil
		}
		defer page.Close()

		if _, err = page.Goto(query); err != nil {
			log.Printf("Could not navigate to %s: %v", query, err)
		}

		// проверка наличия блока "не найдено"
		notFound, err := page.Locator("section.ui-card.w-full").IsVisible()
		if err != nil {
			log.Printf("Error checking visibility: %v", err)
			return nil
		}
		if notFound {
			log.Printf("Event %s not found in city %s", name, city)
			return nil
		}

		eventBlocks, err := page.Locator(".recommendation-item.compilation-tile").All()
		if err != nil {
			log.Printf("Error locating event blocks: %v", err)
			return nil
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		events := []entity.Event{}

		for i, eventBlock := range eventBlocks {
			if len(eventBlocks) > 5 && i == 5 {
				break
			}
			wg.Add(1)
			go func(eventBlock playwright.Locator) {
				defer wg.Done()

				officialName, _ := eventBlock.Locator(".recommendation-item_title.compilation-tile__title").TextContent()
				location, _ := eventBlock.Locator(".recommendation-item_venue.compilation-tile__venue").TextContent()
				ticketsPageUrl, _ := eventBlock.Locator("a").First().GetAttribute("href")

				log.Printf("Parsing event: name-%s, location-%s, url-%s", officialName, location, ticketsPageUrl)

				mu.Lock()
				tickets, err := parseTicketsKassir(ctx, page, cityToUrl[city]+ticketsPageUrl, location)
				if err != nil {
					log.Printf("Error parsing tickets for %s: %v", officialName, err)
				}

				events = append(events, entity.Event{
					Name:    officialName,
					City:    city,
					Tickets: tickets,
				})
				mu.Unlock()
			}(eventBlock)
		}

		wg.Wait()
		return events
	}
}

func parseTicketsKassir(ctx context.Context, page playwright.Page, url, location string) ([]entity.Ticket, error) {
	if _, err := page.Goto(url); err != nil {
		return nil, err
	}

	isOneDate, err := page.Locator("a.focus-ringable.focus-ringable--offset.ui-button--default.ui-button--default-padding.rounded.items-center").Last().IsHidden()
	if err != nil {
		return nil, err
	}

	var tickets []entity.Ticket

	if isOneDate {
		ticket, err := parseSingleTicketKassir(page, location)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	} else {
		datesButtons, err := page.Locator(".event-date-selector-tab").All()
		if err != nil {
			return nil, err
		}

		i := 0
		for _, dateButton := range datesButtons {
			i++
			log.Println("ticket num: ", i)
			if err := dateButton.Click(); err != nil {
				log.Printf("Error clicking date button: %v", err)
				continue
			}

			ticket, err := parseSingleTicketKassir(page, location)
			if err != nil {
				log.Printf("Error parsing ticket: %v", err)
				continue
			}
			tickets = append(tickets, ticket)
		}
	}

	return tickets, nil
}

func parseSingleTicketKassir(page playwright.Page, location string) (entity.Ticket, error) {
	price, err := page.Locator(".flex.cursor-pointer.select-none").First().TextContent()
	if err != nil {
		return entity.Ticket{}, err
	}

	formattedPrice, err := parsePriceKassir(price)
	if err != nil {
		return entity.Ticket{}, err
	}

	date, err := page.Locator(".inline-block span > span:first-child span").TextContent()
	if err != nil {
		return entity.Ticket{}, err
	}

	time, err := page.Locator(".inline-block span > span:last-child span:first-child").First().TextContent()
	if err != nil {
		return entity.Ticket{}, err
	}

	day, err := page.Locator(".inline-block span > span:last-child span span").TextContent()
	if err != nil {
		return entity.Ticket{}, err
	}

	return entity.Ticket{
		Price:    formattedPrice,
		Location: location,
		Date:     date + ", " + day,
		Time:     time,
		Url:      page.URL(),
	}, nil
}

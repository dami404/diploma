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

func parsePriceBileter(priceStr string) (int, error) {
	cleaned := strings.ReplaceAll(priceStr, "от", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	return strconv.Atoi(cleaned)
}

func parseDateTimeBileter(dateTime string) (string, string) {
	cleaned := strings.TrimSpace(dateTime)
	results := strings.Split(cleaned, ",")
	return results[0], results[1]
}

func (r *HTTPDBRepository) ParseBileter(ctx context.Context, city string, name string) []entity.Event {
	select {
	case <-ctx.Done():
		log.Println("Repository.ProfitEvents.ParseBileter: timeout")
		return nil
	default:
		cityToUrl := map[string]string{
			"msk": "https://msk.bileter.ru",
			"spb": "https://www.bileter.ru",
		}
		query := cityToUrl[city] + "/afisha/search?search=" + name

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
		notFound, err := page.Locator(".empty-search").IsVisible()
		if err != nil {
			log.Printf("Error checking visibility: %v", err)
			return nil
		}
		if notFound {
			log.Printf("Event %s not found in city %s", name, city)
			return nil
		}

		eventBlocks, err := page.Locator("div.info-block").All()
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

				officialName, err := eventBlock.Locator("div.name > a").TextContent()
				if err != nil {
					log.Printf("Error parsing officialName: %v", err)
				}
				log.Println("officialName:", officialName)

				location, err := eventBlock.Locator("div.place > a").TextContent()
				log.Println("location:", location)
				if err != nil {
					log.Printf("Error parsing location: %v", err)
				}

				url, err := eventBlock.Locator("div.name > a").First().GetAttribute("href")
				log.Println("url:", url)
				if err != nil {
					log.Printf("Error parsing url: %v", err)
				}

				price, err := eventBlock.Locator("div.price > a > span").TextContent()
				log.Println("price:", price)
				if err != nil {
					log.Printf("Error parsing price: %v", err)
				}

				formattedPrice, err := parsePriceBileter(price)
				log.Println("formattedPrice:", formattedPrice)
				if err != nil {
					log.Printf("Error parsing formattedPrice: %v", err)
				}

				dateTime, _ := eventBlock.Locator("div.date").TextContent()
				log.Println("dateTime:", dateTime)
				if err != nil {
					log.Printf("Error parsing dateTime: %v", err)
				}

				date, time := parseDateTimeBileter(dateTime)
				log.Println("date:", date)
				log.Println("time:", time)
				if err != nil {
					log.Printf("Error parsing date & time: %v", err)
				}
				mu.Lock()

				events = append(events, entity.Event{
					Name: officialName,
					City: city,
					Tickets: []entity.Ticket{
						entity.Ticket{
							Price:    formattedPrice,
							Location: location,
							Date:     date,
							Time:     time,
							Url:      url,
						},
					},
				})
				log.Println(events)
				mu.Unlock()

			}(eventBlock)
		}
		wg.Wait()
		return events
	}
}

func parseTicketsBileter(ctx context.Context, page playwright.Page, url, location string) ([]entity.Ticket, error) {
	if _, err := page.Goto(url); err != nil {
		return nil, err
	}
	return nil, nil
}

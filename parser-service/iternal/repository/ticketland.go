package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/dami404/diploma-parser/iternal/entity"
	"github.com/playwright-community/playwright-go"
)

func parsePriceTicketLand(priceStr string) (int, error) {
	cleaned := strings.ReplaceAll(priceStr, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\u00a0", "")
	cleaned = strings.ReplaceAll(cleaned, "от", "")
	cleaned = strings.ReplaceAll(cleaned, "₽", "")
	cleaned = strings.TrimSpace(cleaned)
	results := strings.Split(cleaned, "–")
	return strconv.Atoi(results[0])
}

func parseDateTicketLand(dateTime string) string {
	cleaned := strings.TrimSpace(dateTime)
	results := strings.Split(cleaned, "•")
	return results[0] + ", " + results[1]
}

func parseTimeTicketLand(time string) string {
	cleaned := strings.ReplaceAll(time, "•", "")
	return cleaned
}

func (r *HTTPDBRepository) ParseTicketLand(ctx context.Context, city string, name string) []entity.Event {
	select {
	case <-ctx.Done():
		log.Println("Repository.ProfitEvents.parseKassirRu: timeout")
		return nil
	default:
		cityToUrl := map[string]string{
			"msk": "https://www.ticketland.ru",
			"spb": "https://spb.ticketland.ru",
			"ekb": "https://ekb.ticketland.ru",
		}
		query := cityToUrl[city] + "/search/performance/?text=" + name
		fmt.Println(query)

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
		notFound, err := page.Locator(".search__body > p.mt-3").IsVisible()
		if err != nil {
			log.Printf("Error checking visibility: %v", err)
			return nil
		}
		if notFound {
			log.Printf("Event %s not found in city %s", name, city)
			return nil
		}

		eventBlocks, err := page.Locator(".card-search.card-search--show").All()
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

				officialName, err := eventBlock.Locator(".card-search__name").TextContent()
				if err != nil {
					log.Printf("Error parsing officialName: %v", err)
				}
				log.Println("officialName:", officialName)

				location, err := eventBlock.Locator(".card-search__building").TextContent()
				log.Println("location:", location)
				if err != nil {
					log.Printf("Error parsing location: %v", err)
				}

				url, err := eventBlock.Locator(".card-search__name").First().GetAttribute("href")
				log.Println("url:", url)
				if err != nil {
					log.Printf("Error parsing url: %v", err)
				}

				price, err := eventBlock.Locator(".card-search__price").TextContent()
				log.Println("price:", price)
				if err != nil {
					log.Printf("Error parsing price: %v", err)
				}

				formattedPrice, err := parsePriceTicketLand(price)
				if err != nil {
					log.Println("error while formatting price:", err)

				}
				log.Println("formattedPrice:", formattedPrice)

				date, _ := eventBlock.Locator(".card-search__date > a").TextContent()
				formattedDate := parseDateTicketLand(date)
				// log.Println("dateTime:", dateTime)
				// if err != nil {
				// 	log.Printf("Error parsing dateTime: %v", err)
				// }
				log.Println("date:", date)

				// date, time := parseDateTimeBileter(dateTime)
				time, _ := eventBlock.Locator(".d-sm-inline").TextContent()
				formattedTime := parseTimeTicketLand(time)
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
							Date:     formattedDate,
							Time:     formattedTime,
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

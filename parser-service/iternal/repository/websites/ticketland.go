package websites

import (
	"context"
	"fmt"
	"log"

	"github.com/dami404/diploma-parser/iternal/entity"
)

func ParseTicketLand(ctx context.Context, city string, name string) []entity.Event {

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
	}
	return nil
}

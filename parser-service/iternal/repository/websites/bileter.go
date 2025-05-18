package websites

import (
	"context"
	"fmt"
	"log"

	"github.com/dami404/diploma-parser/iternal/entity"
)

func ParseBileter(ctx context.Context, city string, name string) []entity.Event {

	select {
	case <-ctx.Done():
		log.Println("Repository.ProfitEvents.parseKassirRu: timeout")
		return nil
	default:
		cityToUrl := map[string]string{
			"msk": "https://msk.bileter.ru",
			"spb": "https://www.bileter.ru",
		}
		query := cityToUrl[city] + "/afisha/search?search=" + name

		fmt.Println(query)
	}
	return nil
}

package service

import (
	"encoding/json"
	"github.com/adlio/trello"
	"github.com/aircjm/cardBox/client"
	"github.com/aircjm/cardBox/client/model"
	"github.com/aircjm/cardBox/config"
	"github.com/aircjm/cardBox/model/request"
	"github.com/aircjm/cardBox/model/response"
	"github.com/aircjm/cardBox/util"
	"log"
)

type AnkiService interface {
}

func AddAnkiNote(cardId string) {
	card, err := client.TrelloCL.GetCard(cardId, trello.Defaults())
	if err != nil {
		panic(err)
	}
	addNoteAnkiRequest := model.AnkiAddNoteRequest{}.GetAnkiAddNote(card)
	response := util.Post(config.AnkiConnect, addNoteAnkiRequest, util.ApplicationJSON)
	log.Println("anki返回的数据是", response)
	ankiResponse := model.AnkiResponse{}
	_ = json.Unmarshal([]byte(response), &ankiResponse)
	log.Println(ankiResponse.Result)
	// 更新anki时间
}

func UpdateAnkiNote(cardId string) {
	card, err := client.TrelloCL.GetCard(cardId, trello.Defaults())
	if err != nil {
		panic(err)
	}
	addNoteAnkiRequest := model.AnkiAddNoteRequest{}.GetAnkiAddNote(card)
	response := util.Post(config.AnkiConnect, addNoteAnkiRequest, util.ApplicationJSON)
	log.Println("anki返回的数据是", response)
	ankiResponse := model.AnkiResponse{}
	_ = json.Unmarshal([]byte(response), &ankiResponse)
	log.Println(ankiResponse.Result)
	// 更新anki时间
}

func TestAnkiConnect() bool {

	response := util.Get(config.AnkiConnect)

	if len(response) != 0 {
		return true
	}

	return false
}

func SaveTrelloToAnkiDecks() {
	boards, err := client.TrelloCL.GetMyBoards(trello.Defaults())

	if err != nil {
		panic(err)
	}

	if !TestAnkiConnect() {
		panic("无法连接anki服务")
	}

	// 需要添加如果已经有deck不要二次添加的逻辑
	for _, board := range boards {
		addDeckRequest := model.AnkiAddDeckRequest{}
		addDeckRequest.Action = "createDeck"
		addDeckRequest.Version = 6
		addDeckRequest.Params.Deck = board.Name
		util.Post(config.AnkiConnect, addDeckRequest, util.ApplicationJSON)
	}
}

// SaveCardToAnki 保存数据到anki
func SaveCardToAnki(Ids []string) {
	for e := range Ids {
		AddAnkiNote(Ids[e])
	}
}

func GetCardList(request request.GetCardListRequest) ([]response.CardResponse, error) {

	var cards []*trello.Card
	cardResponseList := []response.CardResponse{}

	if len(request.BoardId) > 0 {

		board, err := client.TrelloCL.GetBoard(request.BoardId, trello.Defaults())
		if err != nil {
			log.Fatalln(err)
		}
		cards, err = board.GetCards(trello.Defaults())
		log.Println("查询入参有boardId")

		for _, card := range cards {
			cardResponse := response.CardResponse{}

			cardResponse.CardInfo.Id = card.ID
			cardResponse.CardInfo.Name = card.Name

			cardResponseList = append(cardResponseList, cardResponse)
		}
	}
	return cardResponseList, nil
}

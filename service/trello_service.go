package service

import (
	"github.com/adlio/trello"
	"github.com/aircjm/cardBox/client"
	"github.com/aircjm/cardBox/dao"
	"log"
)

// GetRecentlyEditedCard 获取最新的卡片记录
func GetRecentlyEditedCard() (cards []*trello.Card, err error) {
	cards, err = client.TrelloCL.SearchCards("edited:week", trello.Defaults())
	return cards, err
}

// 获取包含对应标签标签的卡片
func GetLabelCard(board trello.Board, labelName string) []*trello.Card {
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	var labelCardLists []*trello.Card
	for _, card := range cards {
		if len(card.Labels) == 0 {
			continue
		}
		for _, label := range card.Labels {
			if label.Name == labelName {
				labelCardLists = append(labelCardLists, card)
			}
		}
	}
	return labelCardLists
}

// GetRecentlyEditedCard 获取最新的卡片记录
func SaveRecentlyEditedCard() {
	cards, err := GetRecentlyEditedCard()
	if err != nil {
		log.Fatalln(err)
	}
	SaveCardsOrm(cards)
}

// SaveCards 批量保存cards 如果有就更新
func SaveCards(cards []*trello.Card) {
	for _, card := range cards {
		go dao.SaveCard(*card)
	}
}

// SaveCards 批量保存cards 如果有就更新
func SaveCardsOrm(cards []*trello.Card) {
	for _, card := range cards {
		dao.SaveCardEntity(*card)
	}
}

func GetBoardAnkiLabelCard(board trello.Board) []*trello.Card {

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	var ankiLabelCardList []*trello.Card

	for _, card := range cards {
		if len(card.Labels) == 0 {
			continue
		}
		for _, label := range card.Labels {
			if label.Name == "anki" {
				ankiLabelCardList = append(ankiLabelCardList, card)
			}
		}
	}
	return ankiLabelCardList
}

// SaveAllCards 保存所有的卡片
func SaveAllCards() {
	boards, err := client.TrelloCL.GetMyBoards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	for _, board := range boards {
		cards, err := board.GetCards(trello.Defaults())
		if err != nil {
			log.Fatal(err)
		}
		go SaveCardsOrm(cards)
	}
}

func GetBoardList() []*trello.Board {
	// 后面需要迁移到查询DB使用，不再直接调用API
	boards, err := client.TrelloCL.GetMyBoards(trello.Defaults())
	if err != nil {
		log.Fatalln(err)
	}
	return boards
}

func ConvertToAnki(list []string) {
	dao.GetCardByCardIdList(list)
	// todo 转换卡片
}

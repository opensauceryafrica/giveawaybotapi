package cron

import (
	"log"

	"github.com/opensaucerers/giveawaybot/repository/v1/giveaway"
	"go.mongodb.org/mongo-driver/bson"
)

// MigrateReplies moves the giveaway replies into the replies collection
func MigrateReplies() {
	// find all giveaways
	// for each giveaway, save the replies ino the replies collection
	// for giveaways := range giveaway.YieldGiveawayByMatch(bson.M{}, config.DefaultPageLimit) {
	// 	log.Printf("Migrating raid history for %d giveaways", len(giveaways))

	// 	for i, giveaway := range giveaways {

	// 		if err := giveaway.ReplyInBatch(); err != nil {
	// 			log.Println(err)
	// 			continue
	// 		}

	// 		log.Printf("len(giveaway.Replies %s): %d", giveaways[i].ID.Hex(), len(giveaways[i].Replies))

	// 		// update raid
	// 		if err := giveaway.ClearReplies(); err != nil {
	// 			log.Println(err)
	// 			continue
	// 		}
	// 	}

	// 	log.Printf("Migrated replies for %d giveaways", len(giveaways))
	// }

	// find all replies
	// for each reply, find the giveaway and add giveaway id to the reply
	for replies := range giveaway.YieldRepliesByMatch(bson.M{}, 61) {
		log.Printf("Updating %d replies", len(replies))

		for _, reply := range replies {

			giveaway := giveaway.Giveaway{
				TweetID: "1654857079423737862",
			}
			err := giveaway.FindByTweet()
			if err != nil {
				log.Println(err)
				continue
			}

			reply.Giveaway = giveaway.ID

			log.Println("reply.Giveaway: ", reply.Giveaway.Hex())

			if err := giveaway.SaveReply(&reply); err != nil {
				log.Println(err)
				continue
			}

			log.Printf("Updated reply %s %s %s", reply.ID, reply.Username, giveaway.ID.Hex())
		}

		log.Printf("Updated %d replies", len(replies))
		return
	}
}

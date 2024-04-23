package graphql

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/graphql-go/graphql"

	"lottery-vote/actors"
	"lottery-vote/utils"
)

// 定义全局变量System来管理Actor System
var System *actor.ActorSystem

// 定义GraphQL Schema
var Schema graphql.Schema

func buildSchema() graphql.Schema {
	voteType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Vote",
		Fields: graphql.Fields{
			"success": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"vote": &graphql.Field{
				Type: voteType,
				Args: graphql.FieldConfigArgument{
					"ticket": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"usernames": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					ticket := params.Args["ticket"].(string)
					usernames := params.Args["usernames"].([]interface{})
					groupedUsernames := make(map[int][]string)

					for _, usernameInterface := range usernames {
						username := usernameInterface.(string)
						index := GetVoteActorIndex(username)
						groupedUsernames[index] = append(groupedUsernames[index], username)
					}

					actorID, _ := strconv.Atoi(strings.Split(ticket, ":")[1])
					ticketVal := strings.Split(ticket, ":")[0]

					validateRes, _ := actors.System.Root.RequestFuture(actors.TicketActorPIDs[actorID], &actors.ValidateTicket{Ticket: ticketVal}, 5*time.Second).Result()
					if validateRes == nil || !validateRes.(*actors.VoteResponse).Success {
						return map[string]interface{}{"success": false}, nil
					}

					var overallSuccess = true
					for index, names := range groupedUsernames {
						voteRequest := &actors.VoteRequest{
							Usernames: names,
						}

						res, _ := actors.System.Root.RequestFuture(actors.VoteActorPIDs[index], voteRequest, 5*time.Second).Result()
						if res == nil {
							overallSuccess = false
							break
						}

						success := res.(*actors.VoteResponse).Success
						if !success {
							overallSuccess = false
							break
						}
					}

					return map[string]interface{}{"success": overallSuccess}, nil
				},
			},
		},
	})

	voteCountType := graphql.NewObject(graphql.ObjectConfig{
		Name: "VoteCount",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"votes": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"getVoteCount": &graphql.Field{
				Type: voteCountType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					username := params.Args["username"].(string)
					var votes int

					fmt.Printf("getVoteCount param : %v\n", username)

					res, err := actors.System.Root.RequestFuture(actors.VoteActorPIDs[GetVoteActorIndex(username)], &actors.VoteQuery{Username: username}, 5*time.Second).Result()
					if err != nil {
						return nil, err
					}
					if res == nil {
						return map[string]interface{}{"username": username, "votes": votes}, nil
					}

					votes = res.(*actors.VoteQueryResponse).Votes
					return map[string]interface{}{"username": username, "votes": votes}, nil
				},
			},
			"getTicket": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					randomActorId := rand.Intn(actors.TicketActorCount)
					resp, err := actors.System.Root.RequestFuture(actors.TicketActorPIDs[randomActorId], &actors.QueryTicket{}, 5*time.Second).Result()
					if err != nil {
						return nil, err
					}
					ticket := fmt.Sprintf("%v:%d", resp.(*actors.TicketResponse).Ticket, randomActorId)
					return ticket, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Mutation: mutationType,
		Query:    queryType,
	})
	if err != nil {
		log.Fatalf("Failed to create new schema, error: %v", err)
	}

	return schema
}

func GetVoteActorIndex(username string) int {
	hash := utils.GetUsernameHash(username)
	return int(hash % uint32(actors.VoteActorCount))
}

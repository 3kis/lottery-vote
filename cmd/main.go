package main

import (
	"lottery-vote/actors"
	"lottery-vote/graphql"
)

func main() {
	actors.InitActorSystem()
	graphql.StartServer()
}

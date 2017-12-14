package main

const gqlSchema = `
schema {
		query: Query
	}

	# The query type, represents all of the entry points into our object graph
	type Query {
		realmByName(name: String!): Realm
		sessionById(id: ID!): Session
		sessionsByRealmId(realmId: ID!): [Session]
		playerById(id: ID!): Player
	}

	type Player {
		id: ID!
		name: String!
		realmId: ID!
	}

	type Realm {
		id: ID!
		name: String!
		title: String
	}

	type Session {
		id: ID!
		realmId: ID!
		name: String
		time: String!
	}

	type PlayerSession {
		playerId: ID!
		sessionId: ID!
		buyin: Int!
		walkout: Int!
	}
`

// const gqlSchema = `
// schema {
// 		query: Query
// 		mutation: Mutation
// 	}
//
// 	# The query type, represents all of the entry points into our object graph
// 	type Query {
// 		realmByName(name: String!): Realm
// 		realmById(id: ID!): Realm
// 		sessionById(id: ID!): Session
// 		sessionsByRealmId(realmId: ID!): [Session]
// 		playerById(id: ID!): Player
// 		playerSessionByPlayerIdSessionId(playerId: ID!, sessionId: ID!): PlayerSession
// 		playerSessionsByPlayerId(playerId: ID!): [PlayerSession]
// 		playerSessionsBySessionId(sessionId: ID!): [PlayerSession]
// 	}
//
// 	# The mutation type, represents all updates we can make to our data
// 	type Mutation {
// 		createRealm(name: String!, title: String): Realm
// 		createPlayer(name: String!, realmId: ID!): Player
// 		createSession(name: String, realmId: ID!, time: String!, playerSessions: [PlayerSession]!)
// 		updateSession(sessionId: ID!, name: String, time: String, playerSessions: [PlayerSession])
// 	}
//
// 	type Player {
// 		id: ID!
// 		name: String!
// 		realmId: ID!
// 		sessions: [Session]
// 		playerSessions: [PlayerSession]
// 		realm: Realm
// 	}
//
// 	type Realm {
// 		id: ID!
// 		name: String!
// 		title: String
// 		sessions: [Session]
// 		players: [Player]
// 		playerSessions: [PlayerSession]
// 	}
//
// 	type Session {
// 		id: ID!
// 		realmId: ID!
// 		name: String
// 		time: String!
// 		players: [Player]
// 		playerSessions: [PlayerSession]
// 	}
//
// 	type PlayerSession {
// 		playerId: ID!
// 		sessionId: ID!
// 		buyin: Int!
// 		walkout: Int!
// 		player: Player
// 		session: Session
// 	}
// `
# Edit-Software-BigGun

# Project:Tablic

Tablic is an implementation of the popular card game "Tablic".Our implementation refers to the two-player version of the game that is played with a standard deck of 52 cards. In this implementation, we have used the Go (Golang) programming language and integrated a PostgreSQL database.

## How to Install and Run the Project

Before you start the Tablic project, you will need following software installed on your machine:

- Go(Golang) : Make sure you have Go installed. If not, you can download and install it from the official website:[Golang Instalation](https://golang.org/)
- PostgreSQL Database: Ensure that you have PostgreSQL installed and running on your system. You can download and install it from the official website: [PostgreSQL Instalation](https://www.postgresql.org/). Remember password you entered in installation process!
  After that you should do the following steps:

* .env file:
  It is necessary to create an .env file with the following constants:

- PORT
- HOST
- DB_PORT
- USER
- PASSWORD
- DB_NAME
- SECRET

* Configure PostgreSQL:
  Create a new PostgreSQL database for the Tablic project.
  Update the database connection details (e.g., username, password, database name) in the .env file.
* Databese migration:
  To create the necessary tables in the database, you need to run the migrate.go script. This script handles the database schema migration and sets up the required tables.
  To run the database migration, open a terminal or command prompt, navigate to the project directory, and execute the following command:

```go
go run migrate.go
```

Running the database migration script is a one-time setup process. You should only run it once, when setting up the application for the first time or when applying database schema changes.

- Run the Project:
  Once the setup is completed, you can run the application by executing the following command in the terminal:

```go
go run main.go
```

## How to Use the Project

The application works through various endpoints. The touchpoints of this communication are considered endpoints. For APIs, an endpoint can include a URL of a server or service. Each endpoint is the location from which APIs can access the resources they need to carry out their function. Endpoints through which this app works are:

### Endpoint for create new deck

This endpoint is called when creating a new deck of cards. The response returned by this endpoint is in JSON format and contains field: "response". The "response" field represents a set of field: "success", "deck_id", "shuffled" and "remaining". In this app, two main fields from this endpoints are "deck_id" which present unique number of deck and "remaining" which present number of cards in deck.

# returns new deck

URL "/cards"
The format of the response received when calling this endpoint is:

```JSON
# json response
{

    "response": {

        "success": true,

        "deck_id": "95v90zxupr95",

        "shuffled": false,

        "remaining": 52

    }

}
```

### Endpoint for create new game

This endpoint is called when user wants to start the new game. User needs to enter his name/nickname and the endpoint is called and the user is placed in a queue where he waits for the other player to join so the game can be started. This endpoint also accepts body parameter "name". "name" body field contain name of player who want to join in game. The response returned by this endpoint is in JSON format and contains fields : "game1", "game2" and "message". The "game1 " field represent set of subfields which present Game model: "ID", "CreatedAt", "UpdatedAt", "DeletedAt", "score", "deckPile","tablePile", "handPile", "collectedPile", "first", "collectedLast", "user_id" and one subfield User which contains information from User model. The "game1" fields represent first created user in game, while field "game2" represent second creted user. The field "game2" consist of the same informaton like field "game1". The field "message" represent a successfully message after created new game.

# returns new game

URL "/addPlayer"
Data that the endpoint requests

```JSON

#json body
{

    "name:":"Player1"

}
```

The format of the response received when calling this endpoint is:

```JSON
# json response
{

    "game1": {

        "ID": 1,

        "CreatedAt": "2023-07-25T12:13:17.5005803+02:00",

        "UpdatedAt": "2023-07-25T12:13:17.5005803+02:00",

        "DeletedAt": null,

        "score": 0,

        "deckPile": "48ri17p220vc",

        "tablePile": "table",

        "handPile": "hand1",

        "collectedPile": "taken1",

        "first": true,

        "collectedLast": false,

        "user_id": 1,

        "User": {

            "ID": 1,

            "CreatedAt": "2023-07-25T12:13:10.408514+02:00",

            "UpdatedAt": "2023-07-25T12:13:10.408514+02:00",

            "DeletedAt": null,

            "name": "Elma"

        }

    },

    "game2": {

        "ID": 2,

        "CreatedAt": "2023-07-25T12:13:17.5129964+02:00",

        "UpdatedAt": "2023-07-25T12:13:17.5129964+02:00",

        "DeletedAt": null,

        "score": 0,

        "deckPile": "48ri17p220vc",

        "tablePile": "table",

        "handPile": "hand2",

        "collectedPile": "taken2",

        "first": false,

        "collectedLast": false,

        "user_id": 2,

        "User": {

            "ID": 2,

            "CreatedAt": "2023-07-25T12:13:16.7659953+02:00",

            "UpdatedAt": "2023-07-25T12:13:16.7659953+02:00",

            "DeletedAt": null,

            "name": "Benjo"

        }

    },

    "message": "Game has started"

}

```

### Endpoint for return cards from hand of player and cards from table

This endpoint is called when player want see his cards from hands and cards from table. This endpoint must be secured from other player, because the other player cannot see his cards. Protection is enabled through the implementation of JWT tokens, where the player, during his game, contains a token with which he is authenticated. Each secure route contains a "Middleware" function that performs certain authentication. The parameters received by this endpoint are "userId" and "deckId", which are used to find player cards in a particular game. The response returned by this endpoint is in JSON format and contains two fields: "Cards_from_table" and "User_hand_cards". "Cards_from_table" field represents a set of cards on the table containing the subfields "image", "value", "suit" and "code", while "User_hand_cards" represents a set of cards of player containing the same subfields like first field.

# returns player hand cards and cards from table

URL "/cards/:userId/:deckId"
The format of the response received when calling this endpoint is:

```JSON
# json response
{

    "Cards_from_table": {

        "cards": [

            {

                "image": "https://deckofcardsapi.com/static/img/8C.png",

                "value": "8",

                "suit": "CLUBS",

                "code": "8C"

            },

            {

                "image": "https://deckofcardsapi.com/static/img/0D.png",

                "value": "10",

                "suit": "DIAMONDS",

                "code": "0D"

            },

            {

                "image": "https://deckofcardsapi.com/static/img/7D.png",

                "value": "7",

                "suit": "DIAMONDS",

                "code": "7D"

            },

            {

                "image": "https://deckofcardsapi.com/static/img/0C.png",

                "value": "10",

                "suit": "CLUBS",

                "code": "0C"

            }

        ],

        "remaining": 4

    },

    "User_hand_cards": [

        {

            "image": "https://deckofcardsapi.com/static/img/9H.png",

            "value": "9",

            "suit": "HEARTS",

            "code": "9H"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/0S.png",

            "value": "10",

            "suit": "SPADES",

            "code": "0S"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/4C.png",

            "value": "4",

            "suit": "CLUBS",

            "code": "4C"

        }

    ]

}
```

### Endpoint for throw card on table

This endpoint is called when player want to throw card on table. This endpoint must be secured from other player, because the other player from other game cannot throw his card. Protection is enabled through the implementation of JWT tokens, where the player, during his game, contains a token with which he is authenticated. Payload in token hold deck_id which game can be authorized. The parameters received by this endpoint are "cardCode", "deckId" and "playerPile", where "cardCode" present card who player want throw, "deckId" present game in which the player plays and unique deck with cards which is used in game, "playerPile" present cards in hands from where the player throws the card. The response returned by this endpoint is in JSON format and contains two fields: "message", "table_cards" and "user_hand_cards". "Message" field represent message after successfully throw card ("The card is thrown on the table"), "table_cards" field represents a set of cards on the table containing the subfields "image", "value", "suit" and "code", while "user_hand_cards" represents a set of cards of player containing the same subfields like first field.

# returns player hand cards and cards from table after throw card on table

URL "/throwCard/:cardCode/:deckId/:playerPile"
The format of the response received when calling this endpoint is:

```JSON
# json response

{

    "message": "The card is thrown on the table",

    "table_cards": [

        {

            "image": "https://deckofcardsapi.com/static/img/2C.png",

            "value": "2",

            "suit": "CLUBS",

            "code": "2C"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/0H.png",

            "value": "10",

            "suit": "HEARTS",

            "code": "0H"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/9D.png",

            "value": "9",

            "suit": "DIAMONDS",

            "code": "9D"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/KH.png",

            "value": "KING",

            "suit": "HEARTS",

            "code": "KH"

        }

    ],

    "user_hand_cards": [

        {

            "image": "https://deckofcardsapi.com/static/img/QH.png",

            "value": "QUEEN",

            "suit": "HEARTS",

            "code": "QH"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/9S.png",

            "value": "9",

            "suit": "SPADES",

            "code": "9S"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/JH.png",

            "value": "JACK",

            "suit": "HEARTS",

            "code": "JH"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/2S.png",

            "value": "2",

            "suit": "SPADES",

            "code": "2S"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/5S.png",

            "value": "5",

            "suit": "SPADES",

            "code": "5S"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/0C.png",

            "value": "10",

            "suit": "CLUBS",

            "code": "0C"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/8H.png",

            "value": "8",

            "suit": "HEARTS",

            "code": "8H"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/JD.png",

            "value": "JACK",

            "suit": "DIAMONDS",

            "code": "JD"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/KS.png",

            "value": "KING",

            "suit": "SPADES",

            "code": "KS"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/QD.png",

            "value": "QUEEN",

            "suit": "DIAMONDS",

            "code": "QD"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/9H.png",

            "value": "9",

            "suit": "HEARTS",

            "code": "9H"

        }

    ]

}
```

### Endpoint for take card from table

This endpoint is called when player want to take card from table. This endpoint must be secured from other player, because the other player from other game cannot throw his card. Protection is enabled through the implementation of JWT tokens, where the player, during his game, contains a token with which he is authenticated. Payload in token hold deck_id which game can be authorized. The parameters received by this endpoint are "deckId", "handPile" and "takenPile", where "deckId" present game in which the player plays and unique deck with cards which is used in game, "handPile" present cards in hands and this parameter use for finding DB result of player and "takenPile" present cards where the player store the card from table. This endpoint also accepts body parameters:"hand_card" and "taken_cards". "hand_card" body field contains sign of card which player want to take other cards/card from table, while "taken_cards" contains signs of group of cards separated by ";" which can be taken. The response returned by this endpoint is in JSON format and contains two fields: "response", "table_cards" and "user_hand_cards". "Message" field represent message after successfully taken card ("Cards are moved from hand and table pile to taken pile"), "table_cards" field represents a set of cards on the table containing the subfields "image", "value", "suit" and "code", while "user_hand_cards" represents a set of cards of player containing the same subfields like first field.

# returns player hand cards and cards from table after take card/cards from table

URL "/takecardsfromtable/:deckId/:handPile/:takenPile"
Data that the endpoint requests

```JSON
#json body
{

    "hand_card" : "AD",

    "taken_cards": "AS;0D,1D;6S,5H"

}

```

The format of the response received when calling this endpoint is:

```JSON
# json response

{

    "response": "Cards are moved from hand and table pile to taken pile",

    "table_cards": [

        {

            "image": "https://deckofcardsapi.com/static/img/2C.png",

            "value": "2",

            "suit": "CLUBS",

            "code": "2C"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/0H.png",

            "value": "10",

            "suit": "HEARTS",

            "code": "0H"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/9D.png",

            "value": "9",

            "suit": "DIAMONDS",

            "code": "9D"

        }

    ],

    "user_hand_cards": [

        {

            "image": "https://deckofcardsapi.com/static/img/0S.png",

            "value": "10",

            "suit": "SPADES",

            "code": "0S"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/7D.png",

            "value": "7",

            "suit": "DIAMONDS",

            "code": "7D"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/AS.png",

            "value": "ACE",

            "suit": "SPADES",

            "code": "AS"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/6C.png",

            "value": "6",

            "suit": "CLUBS",

            "code": "6C"

        },

        {

            "image": "https://deckofcardsapi.com/static/img/2H.png",

            "value": "2",

            "suit": "HEARTS",

            "code": "2H"

        }

    ]

}
```

### Endpoint for create authorization JWT token

This endpoint is used for create JWT authorization token. Token should be created after user get in game. JWT token used for authorization routes: route for listing player cards in hand, route for throwing card on table and route for take card/cards from table. The parameters received by this endpoint are "userID" and "deckId", where "deckId" present game in which the player plays and unique deck with cards which is used in game, "userId" present id of user. The response returned by this endpoint is in JSON format and contain tfield: "token". The "token" field represent JWT token after successfully created them. JWT token payload consists of fields "user_id", "deck_id" and "exp". Field "user_id" uses for checing authorization in route for listing player hand cards, while "deck_id" uses for authorization in route for throwing card on table and takeing card/cards from table.

# returns JWT player token

URL "/gettoken/:userId/:deckId"
The format of the response received when calling this endpoint is:

```JSON
# json response
{

    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZWNrX2lkIjoiNTVpeHk3enA5eXd1IiwiZXhwIjoxNjkxNjg0NTY2LCJ1c2VyX2lkIjoiMiJ9.a87eqrvtIkj3Z0fZ156rQbfB__w2nIKNbgVaCzMkSME"

}
```

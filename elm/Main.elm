module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (id, class, style)
import Html.Events exposing (onClick)
import WebSocket
import Json.Decode exposing (..)
import Robots exposing (robot)
import Keyboard


main =
    Html.programWithFlags
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }



-- MODEL


type alias Flags =
    { socketURL : String
    , currentPlayerName : String
    , winningCount : Int
    }


type Role
    = Admin
    | Normal


type alias Player =
    { name : String
    , count : Int
    , role : Role
    }


type GameState
    = WaitingForPlayers
    | WaitingForCountdown
    | CountingDown
    | Started
    | Finished


type alias Game =
    { key : String
    , state : GameState
    , players : List Player
    }


type alias Model =
    { socketURL : String
    , currentPlayerName : String
    , countdown : String
    , game : Game
    , winningCount : Int
    }


init : Flags -> ( Model, Cmd Msg )
init flags =
    ( { socketURL = flags.socketURL
      , currentPlayerName = flags.currentPlayerName
      , countdown = ""
      , winningCount = flags.winningCount
      , game =
            { key = ""
            , state = WaitingForPlayers
            , players = []
            }
      }
    , Cmd.none
    )



-- util function


isAdmin model =
    let
        currentPlayer =
            List.filter
                (\p -> p.name == model.currentPlayerName)
                model.game.players
                |> List.head
    in
        case currentPlayer of
            Just p ->
                p.role == Admin

            Nothing ->
                False



-- JSON decoders


type ServerMsg
    = Countdown String
    | State Game


decodeSocketMsg : Decoder ServerMsg
decodeSocketMsg =
    (field "msgType" string) |> andThen decodeSocketMsgBody


decodeSocketMsgBody : String -> Decoder ServerMsg
decodeSocketMsgBody msgType =
    case msgType of
        "countdown" ->
            Json.Decode.map Countdown (field "tick" string)

        "state" ->
            Json.Decode.map State (field "game" decodeGame)

        _ ->
            fail "expecting some kind of point"


decodeGame : Decoder Game
decodeGame =
    map3 Game
        (field "key" string)
        ((field "state" string) |> andThen decodeGameState)
        (field "players" (list decodePlayer))


decodeGameState : String -> Decoder GameState
decodeGameState state =
    succeed <|
        case state of
            "WaitingForCountdown" ->
                WaitingForCountdown

            "CountingDown" ->
                CountingDown

            "Started" ->
                Started

            "Finished" ->
                Finished

            _ ->
                WaitingForPlayers


decodePlayer : Decoder Player
decodePlayer =
    map3 Player
        (field "name" string)
        (field "count" int)
        ((field "role" string) |> andThen decodeRole)


decodeRole : String -> Decoder Role
decodeRole role =
    succeed <|
        case role of
            "Admin" ->
                Admin

            _ ->
                Normal



-- UPDATE


type Msg
    = CloseToNewPlayers
    | Start
    | Increment
    | Reset
    | MsgFromSocket String
    | NoOp


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        CloseToNewPlayers ->
            ( model, WebSocket.send model.socketURL "CloseToNewPlayers" )

        Start ->
            ( model, WebSocket.send model.socketURL "Start" )

        Increment ->
            ( model, WebSocket.send model.socketURL "Increment" )

        Reset ->
            ( model, WebSocket.send model.socketURL "Reset" )

        MsgFromSocket socketMsg ->
            case Json.Decode.decodeString decodeSocketMsg socketMsg of
                Ok val ->
                    case val of
                        Countdown tick ->
                            ( { model | countdown = tick }, Cmd.none )

                        State game ->
                            ( { model | game = game }, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        NoOp ->
            ( model, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ WebSocket.listen model.socketURL MsgFromSocket
        , Keyboard.ups (onKeyDown model)
        ]


onKeyDown model key =
    case key of
        32 ->
            Increment

        27 ->
            Reset

        _ ->
            NoOp



-- VIEW


view : Model -> Html Msg
view model =
    let
        foo =
            viewPlayer model
    in
        div [ id "racetrack", onClick <| inc model ]
            [ dialog model
            , div [] (List.indexedMap foo model.game.players)
            ]


inc model =
    case model.game.state of
        Started ->
            Increment

        _ ->
            NoOp


dialog model =
    case model.game.state of
        WaitingForPlayers ->
            let
                controls =
                    if isAdmin model && (List.length model.game.players) >= 2 then
                        button [ onClick CloseToNewPlayers ] [ text "Close to new players" ]
                    else
                        text ""
            in
                div [ id "dialog" ]
                    [ h1 [] [ text "Waiting for players to join" ]
                    , div [] [ controls ]
                    , p [] [ text "Invite players by sharing URL" ]
                    , p [] [ text "Min two players" ]
                    , p [] [ text "Max five players" ]
                    ]

        WaitingForCountdown ->
            let
                controls =
                    if isAdmin model then
                        button [ onClick Start ] [ text "Start countdown" ]
                    else
                        text ""
            in
                div [ id "dialog" ]
                    [ h1 [] [ text "Get ready" ]
                    , div []
                        [ controls ]
                    , p [] [ text "Press spacebar or tap screen on mobile to jump" ]
                    ]

        CountingDown ->
            div [ id "dialog", class "get-set" ]
                [ h1 [] [ text model.countdown ]
                ]

        Finished ->
            let
                leaderboard =
                    List.sortBy (\p -> p.count) model.game.players
                        |> List.reverse

                controls =
                    if isAdmin model then
                        button [ onClick Reset ] [ text "Play again" ]
                    else
                        text ""
            in
                case List.head leaderboard of
                    Just winner ->
                        div [ id "dialog" ]
                            [ h1 [] [ text (winner.name ++ " wins!") ]
                            , div [] [ controls ]
                            ]

                    Nothing ->
                        text ""

        _ ->
            text ""


viewPlayer : Model -> Int -> Player -> Html msg
viewPlayer model index player =
    let
        -- Work out vertical placement of player
        playerHeight =
            10

        numIntervalsToWin =
            toFloat model.winningCount

        viewHeight =
            100

        raceHeight =
            viewHeight - playerHeight

        intervalHeight =
            raceHeight / numIntervalsToWin

        height =
            intervalHeight * (toFloat player.count)

        -- Work out horizontal placement of player
        numPlayers =
            List.length model.game.players

        viewWidth =
            100

        intervalWidth =
            viewWidth / toFloat (numPlayers + 1)

        width =
            intervalWidth * (toFloat <| index + 1)

        color =
            case index of
                1 ->
                    -- Green
                    "rgb(42, 166, 61)"

                2 ->
                    -- Blue
                    "rgb(3, 124, 210)"

                3 ->
                    -- Purple
                    "rgb(166, 123, 248)"

                4 ->
                    -- Orange
                    "#FF8000"

                _ ->
                    -- Red
                    -- Is used for the first player.
                    "rgb(239, 25, 77)"
    in
        robot color height width player.name

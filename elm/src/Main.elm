module Main exposing (..)

import Array
import Bootstrap.Alert
import Bootstrap.Button as Button
import Bootstrap.ButtonGroup as ButtonGroup
import Bootstrap.CDN
import Bootstrap.Card as Card
import Bootstrap.Card.Block
import Bootstrap.Form.Input as Input
import Bootstrap.Grid as Grid
import Bootstrap.Grid.Col
import Bootstrap.Grid.Row
import Bootstrap.Navbar as Navbar
import Bootstrap.Text as Text
import Bootstrap.Utilities.Spacing as Spacing
import Browser
import Delay
import Html exposing (Html, div, text)
import Html.Attributes
import Http
import Json.Decode as D



--MAIN


main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }



--MODEL


type LoginRequestStatusModel
    = InProgress
    | CredentialsFailure String
    | ConnectionFailure
    | Success


type alias RegisterRequestStatus =
    LoginRequestStatusModel


type alias LoginPageModel =
    { username : Username
    , password : String
    , status : LoginRequestStatusModel
    }


type alias RegisterPageModel =
    { username : Username
    , password : String
    , status : RegisterRequestStatus
    }


type alias LobbyPageModel =
    { inputRoom : String
    , errorMessage : Maybe String
    }


type alias GameResult =
    { isOver : Bool
    , isDraw : Bool
    , winner : Username
    }


type alias GameModel =
    { p1mancala : Int
    , p2mancala : Int
    , p1cells : List Int
    , p2cells : List Int
    , player1 : String
    , player2 : String
    , p1score : Int
    , p2score : Int
    , result : GameResult
    , next_player : Int
    , errorMessage : Maybe String
    }


type alias Username =
    String


type MainWindow
    = LoginPage
    | SignUpPage
    | LobbySelectionPage
    | InGamePage


type LobbyRole
    = Host
    | Guest
    | None


type alias Model =
    { page : MainWindow
    , loginData : LoginPageModel
    , registerData : RegisterPageModel
    , username : Maybe Username
    , currentRoom : Maybe String
    , navbarState : Navbar.State
    , lobbyData : LobbyPageModel
    , gameData : GameModel
    , lobbyRole : LobbyRole
    , cyclingUpdate : Bool
    }


init : () -> ( Model, Cmd Msg )
init _ =
    let
        ( navbarState, navbarCmd ) =
            Navbar.initialState NavbarMsg
    in
    ( { page = LoginPage
      , loginData = LoginPageModel "" "" InProgress
      , registerData = LoginPageModel "" "" InProgress
      , username = Nothing
      , navbarState = navbarState
      , lobbyData = LobbyPageModel "" Nothing
      , currentRoom = Nothing
      , gameData = GameModel 0 0 [ 4, 4, 4, 4, 4, 4 ] [ 4, 4, 4, 4, 4, 4 ] "A" "B" 0 0 (GameResult False False "") 1 Nothing
      , lobbyRole = None
      , cyclingUpdate = False
      }
    , navbarCmd
    )



--UPDATE


type LoginMessages
    = GotLoginResp (Result Http.Error String)
    | LoginPressed
    | UsernameInput String
    | PasswordInput String


type RegisterMessages
    = GotRegisterResp (Result Http.Error String)
    | RegisterPressed
    | RegUsernameInput String
    | RegPasswordInput String


type LobbyMsg
    = HostGame
    | JoinGame String
    | JoinInput String
    | GotHostResponse (Result Http.Error String)
    | GotJoinResponse (Result Http.Error String)


type GameMsg
    = ExitGame
    | MakeTurn Int
    | GotStateResponse (Result Http.Error String)
    | GotTurnResponse (Result Http.Error String)
    | Refresh


type Msg
    = LoginMessages LoginMessages
    | RegisterMessages RegisterMessages
    | NavbarMsg Navbar.State
    | SelectPageMsg MainWindow
    | LobbyMsg LobbyMsg
    | GameMsg GameMsg
    | LogOut


decodeAuthResponse str =
    case D.decodeString (D.maybe (D.field "error" D.string)) str of
        Ok message ->
            case message of
                Just string ->
                    CredentialsFailure string

                Nothing ->
                    CredentialsFailure "Unknown error"

        Err _ ->
            Success


updateLogin : LoginMessages -> Model -> ( Model, Cmd Msg )
updateLogin msg model =
    case msg of
        GotLoginResp result ->
            case result of
                Ok message ->
                    case decodeAuthResponse message of
                        Success ->
                            let
                                loginDataRec =
                                    model.loginData
                            in
                            ( { model | loginData = { loginDataRec | status = Success }, username = Just loginDataRec.username }, Delay.after 2 Delay.Second (SelectPageMsg LobbySelectionPage) )

                        CredentialsFailure str ->
                            let
                                loginDataRec =
                                    model.loginData
                            in
                            ( { model | loginData = { loginDataRec | status = CredentialsFailure str } }, Cmd.none )

                        InProgress ->
                            ( model, Cmd.none )

                        ConnectionFailure ->
                            ( model, Cmd.none )

                Err _ ->
                    let
                        loginDataRec =
                            model.loginData
                    in
                    ( { model | loginData = { loginDataRec | status = ConnectionFailure } }, Cmd.none )

        LoginPressed ->
            ( model
            , Http.post
                { url = "/api/user/login"
                , body = Http.stringBody "application/x-www-form-urlencoded" ("username=" ++ model.loginData.username ++ "&password=" ++ model.loginData.password)
                , expect = Http.expectString (\a -> LoginMessages (GotLoginResp a))
                }
            )

        UsernameInput string ->
            let
                loginDataRec =
                    model.loginData
            in
            ( { model | loginData = { loginDataRec | username = string } }, Cmd.none )

        PasswordInput string ->
            let
                loginDataRec =
                    model.loginData
            in
            ( { model | loginData = { loginDataRec | password = string } }, Cmd.none )


updateRegister : RegisterMessages -> Model -> ( Model, Cmd Msg )
updateRegister msg model =
    case msg of
        GotRegisterResp result ->
            case result of
                Ok message ->
                    case decodeAuthResponse message of
                        Success ->
                            let
                                registerDataRec =
                                    model.registerData
                            in
                            ( { model | registerData = { registerDataRec | status = Success }, username = Just registerDataRec.username }, Delay.after 2 Delay.Second (SelectPageMsg LobbySelectionPage) )

                        CredentialsFailure str ->
                            let
                                registerDataRec =
                                    model.registerData
                            in
                            ( { model | registerData = { registerDataRec | status = CredentialsFailure str } }, Cmd.none )

                        InProgress ->
                            ( model, Cmd.none )

                        ConnectionFailure ->
                            ( model, Cmd.none )

                Err _ ->
                    let
                        registerDataRec =
                            model.registerData
                    in
                    ( { model | registerData = { registerDataRec | status = ConnectionFailure } }, Cmd.none )

        RegisterPressed ->
            ( model
            , Http.post
                { url = "/api/user/new"
                , body = Http.stringBody "application/x-www-form-urlencoded" ("username=" ++ model.registerData.username ++ "&password=" ++ model.registerData.password)
                , expect = Http.expectString (\a -> RegisterMessages (GotRegisterResp a))
                }
            )

        RegUsernameInput string ->
            let
                registerDataRec =
                    model.registerData
            in
            ( { model | registerData = { registerDataRec | username = string } }, Cmd.none )

        RegPasswordInput string ->
            let
                registerDataRec =
                    model.registerData
            in
            ( { model | registerData = { registerDataRec | password = string } }, Cmd.none )


processJoinLobbyResponse : String -> Model -> ( Model, Cmd Msg )
processJoinLobbyResponse str model =
    case D.decodeString (D.maybe (D.field "error" D.string)) str of
        Ok message ->
            case message of
                Just string ->
                    let
                        ldr =
                            model.lobbyData
                    in
                    ( { model | lobbyData = { ldr | errorMessage = Just string } }, Cmd.none )

                Nothing ->
                    let
                        ldr =
                            model.lobbyData
                    in
                    ( { model | lobbyData = { ldr | errorMessage = Just "Unknown error" } }, Cmd.none )

        Err _ ->
            ( { model | currentRoom = Just model.lobbyData.inputRoom, lobbyRole = Guest }, Delay.after 10 Delay.Millisecond (SelectPageMsg InGamePage) )


updateLobby : LobbyMsg -> Model -> ( Model, Cmd Msg )
updateLobby msg model =
    case msg of
        HostGame ->
            ( model
            , Http.post
                { url = "/api/lobby/create"
                , body = Http.stringBody "application/x-www-form-urlencoded" ""
                , expect = Http.expectString (\a -> LobbyMsg (GotHostResponse a))
                }
            )

        JoinGame _ ->
            ( model
            , Http.post
                { url = "/api/lobby/join"
                , body = Http.stringBody "application/x-www-form-urlencoded" ("room=" ++ model.lobbyData.inputRoom)
                , expect = Http.expectString (\a -> LobbyMsg (GotJoinResponse a))
                }
            )

        GotHostResponse result ->
            case result of
                Ok value ->
                    ( { model | currentRoom = Just value, lobbyRole = Host }, Delay.after 10 Delay.Millisecond (SelectPageMsg InGamePage) )

                Err _ ->
                    let
                        lobbyDataRec =
                            model.lobbyData
                    in
                    ( { model | lobbyData = { lobbyDataRec | errorMessage = Just "Connection error" } }, Cmd.none )

        GotJoinResponse result ->
            case result of
                Ok value ->
                    processJoinLobbyResponse value model

                Err _ ->
                    let
                        lobbyDataRec =
                            model.lobbyData
                    in
                    ( { model | lobbyData = { lobbyDataRec | errorMessage = Just "Connection error" } }, Cmd.none )

        JoinInput string ->
            let
                ldr =
                    model.lobbyData
            in
            ( { model | lobbyData = { ldr | inputRoom = string } }, Cmd.none )


type alias McCell =
    { state : Int }


cell : D.Decoder McCell
cell =
    D.map McCell (D.field "score" D.int)


gameres : D.Decoder GameResult
gameres =
    D.map3 GameResult
        (D.field "game_over" D.bool)
        (D.field "is_draw" D.bool)
        (D.field "winner" D.string)


type alias McScores =
    { player1_score : Int
    , player2_score : Int
    , player1_cells : List McCell
    , player2_cells : List McCell
    , player1_mancala : McCell
    , player2_mancala : McCell
    }


type alias McMetadata =
    { result : GameResult
    , players : List String
    , next_player : Int
    }


scoresBoard : D.Decoder McScores
scoresBoard =
    D.map6 McScores
        (D.field "player1_score" D.int)
        (D.field "player2_score" D.int)
        (D.field "player1_cells" (D.list cell))
        (D.field "player2_cells" (D.list cell))
        (D.field "player1_mancala" cell)
        (D.field "player2_mancala" cell)


metadataBoard : D.Decoder McMetadata
metadataBoard =
    D.map3 McMetadata
        (D.field "result" gameres)
        (D.field "players" (D.list D.string))
        (D.field "next_player" D.int)


getGameStateRequest =
    Http.get
        { url = "/api/lobby/state"
        , expect = Http.expectString (\a -> GameMsg (GotStateResponse a))
        }


processStateJSON : String -> Model -> Cmd Msg -> ( Model, Cmd Msg )
processStateJSON str model cmd =
    let
        scores =
            D.decodeString scoresBoard str
    in
    let
        meta =
            D.decodeString metadataBoard str
    in
    let
        gdr =
            model.gameData
    in
    case scores of
        Ok newScores ->
            case meta of
                Ok newMeta ->
                    let
                        players =
                            Array.fromList newMeta.players
                    in
                    ( { model
                        | gameData =
                            { gdr
                                | p1score = newScores.player1_score
                                , p2score = newScores.player2_score
                                , p1mancala = newScores.player1_mancala.state
                                , p2mancala = newScores.player2_mancala.state
                                , p1cells = List.map (\a -> a.state) newScores.player1_cells
                                , p2cells = List.map (\a -> a.state) newScores.player2_cells
                                , errorMessage = Nothing
                                , result = newMeta.result
                                , next_player = newMeta.next_player
                                , player1 =
                                    case Array.get 0 players of
                                        Just name ->
                                            name

                                        Nothing ->
                                            ""
                                , player2 =
                                    case Array.get 1 players of
                                        Just name ->
                                            name

                                        Nothing ->
                                            ""
                            }
                      }
                    , cmd
                    )

                Err _ ->
                    ( { model | gameData = { gdr | errorMessage = Just "Invalid message received" } }, cmd )

        Err _ ->
            ( { model | gameData = { gdr | errorMessage = Just "Invalid message received" } }, cmd )


processGameResponse : String -> Bool -> Model -> ( Model, Cmd Msg )
processGameResponse response getStateNext model =
    let
        nextCommand =
            if getStateNext then
                Delay.after 1 Delay.Second (GameMsg Refresh)

            else
                Cmd.none
    in
    let
        gdr =
            model.gameData
    in
    case D.decodeString (D.field "error" (D.maybe D.string)) response of
        Ok message ->
            case message of
                Just string ->
                    ( { model | gameData = { gdr | errorMessage = Just string } }, nextCommand )

                Nothing ->
                    ( { model | gameData = { gdr | errorMessage = Just "Unknown error" } }, nextCommand )

        Err _ ->
            processStateJSON response model nextCommand


updateGame : GameMsg -> Model -> ( Model, Cmd Msg )
updateGame msg model =
    case msg of
        ExitGame ->
            ( { model | gameData = GameModel 0 0 [] [] "" "" 0 0 (GameResult False False "") 1 Nothing }, Delay.after 1 Delay.Millisecond (SelectPageMsg LobbySelectionPage) )

        MakeTurn int ->
            ( model
            , Http.post
                { url = "/api/lobby/turn"
                , body = Http.stringBody "application/x-www-form-urlencoded" ("cell=" ++ String.fromInt int)
                , expect = Http.expectString (\a -> GameMsg (GotTurnResponse a))
                }
            )

        GotTurnResponse result ->
            case result of
                Ok value ->
                    processGameResponse value False model

                Err _ ->
                    let
                        gdr =
                            model.gameData
                    in
                    ( { model | gameData = { gdr | errorMessage = Just "Connection error" } }, Cmd.none )

        GotStateResponse result ->
            case result of
                Ok value ->
                    processGameResponse value model.cyclingUpdate model

                Err _ ->
                    let
                        gdr =
                            model.gameData
                    in
                    ( { model | gameData = { gdr | errorMessage = Just "Connection error" } }, Delay.after 1 Delay.Second (GameMsg Refresh) )

        Refresh ->
            ( model, getGameStateRequest )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case Debug.log "msg" msg of
        LoginMessages lMsg ->
            updateLogin lMsg model

        RegisterMessages rMsg ->
            updateRegister rMsg model

        NavbarMsg state ->
            ( { model | navbarState = state }, Cmd.none )

        SelectPageMsg spm ->
            if spm == InGamePage then
                if model.cyclingUpdate then
                    ( { model | page = spm }, Cmd.none )

                else
                    ( { model | page = spm, cyclingUpdate = True }, getGameStateRequest )

            else
                ( { model | page = spm, cyclingUpdate = False }, Cmd.none )

        LogOut ->
            init ()

        LobbyMsg lobbyMsg ->
            updateLobby lobbyMsg model

        GameMsg gameMsg ->
            updateGame gameMsg model



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Navbar.subscriptions model.navbarState NavbarMsg



-- VIEW
--maxWidth100 : Html.Attribute


maxWidth100 =
    Html.Attributes.style "max-width" "100%"


view : Model -> Html Msg
view model =
    preprocessView (Debug.log "model" model)


preprocessView : Model -> Html Msg
preprocessView model =
    Grid.container [] [ Bootstrap.CDN.stylesheet, navBar model, processView model ]


navBar : Model -> Html Msg
navBar model =
    case model.username of
        Nothing ->
            Navbar.config NavbarMsg
                |> Navbar.withAnimation
                |> Navbar.success
                |> Navbar.brand [] [ text "Mancala" ]
                |> Navbar.items
                    [ Navbar.itemLink [] [ Button.button [ Button.onClick (SelectPageMsg LoginPage) ] [ text "Log in" ] ]
                    , Navbar.itemLink [] [ Button.button [ Button.onClick (SelectPageMsg SignUpPage) ] [ text "Register" ] ]
                    ]
                |> Navbar.view model.navbarState

        Just username ->
            Navbar.config NavbarMsg
                |> Navbar.withAnimation
                |> Navbar.success
                |> Navbar.brand [] [ text "Mancala" ]
                |> Navbar.items
                    (if model.page == InGamePage then
                        [ Navbar.itemLink [] [ Button.button [] [ text username ] ]
                        , Navbar.itemLink [] [ Button.button [ Button.onClick (GameMsg ExitGame) ] [ text "Exit game" ] ]
                        , Navbar.itemLink [] [ Button.button [ Button.onClick LogOut ] [ text "Log out" ] ]
                        ]

                     else
                        [ Navbar.itemLink [] [ Button.button [] [ text username ] ]
                        , Navbar.itemLink [] [ Button.button [ Button.onClick LogOut ] [ text "Log out" ] ]
                        ]
                    )
                |> Navbar.view model.navbarState


processView : Model -> Html Msg
processView model =
    case model.page of
        LoginPage ->
            loginForm model

        SignUpPage ->
            signupForm model

        LobbySelectionPage ->
            lobbyView model

        InGamePage ->
            div []
                [ case model.currentRoom of
                    Just id ->
                        Html.h4 [] [ text id ]

                    Nothing ->
                        Html.h4 [] [ text "Not in any room???" ]
                , gameView model
                ]


lobbyView : Model -> Html Msg
lobbyView model =
    div []
        [ Grid.row [ Bootstrap.Grid.Row.attrs [ maxWidth100 ] ]
            [ Grid.col [ Bootstrap.Grid.Col.textAlign Text.alignXsCenter ]
                [ Button.button [ Button.onClick (LobbyMsg HostGame), Button.attrs [ Spacing.m5 ], Button.outlinePrimary ] [ Html.h4 [] [ text "Host game" ] ]
                ]
            , Grid.col [ Bootstrap.Grid.Col.textAlign Text.alignXsCenter ]
                [ Input.text [ Input.onInput (\a -> LobbyMsg (JoinInput a)), Input.placeholder "Room ID", Input.attrs [ Spacing.m5, maxWidth100 ] ]
                , Button.button [ Button.onClick (LobbyMsg (JoinGame model.lobbyData.inputRoom)), Button.outlinePrimary ] [ Html.h4 [] [ text "Join game" ] ]
                ]
            ]
        , Grid.row [ Bootstrap.Grid.Row.attrs [ maxWidth100 ] ]
            [ Grid.col []
                [ lobbyAlert model
                ]
            ]
        ]


lobbyAlert : Model -> Html Msg
lobbyAlert model =
    case model.lobbyData.errorMessage of
        Just string ->
            Bootstrap.Alert.simpleWarning [ Spacing.m5 ] [ text string ]

        Nothing ->
            div [] []


loginForm : Model -> Html Msg
loginForm model =
    Card.config [ Card.outlinePrimary, Card.attrs [ Spacing.m5 ] ]
        |> Card.headerH1 [] [ text "Log in to proceed" ]
        |> Card.block []
            [ Bootstrap.Card.Block.text []
                [ Input.text [ Input.id "usernameInput", Input.onInput (\a -> LoginMessages (UsernameInput a)), Input.value model.loginData.username, Input.attrs [ Spacing.m3, maxWidth100 ] ]
                , Input.password [ Input.id "passwordInput", Input.onInput (\a -> LoginMessages (PasswordInput a)), Input.value model.loginData.password, Input.attrs [ Spacing.m3, maxWidth100 ] ]
                , Button.button [ Button.outlinePrimary, Button.onClick (LoginMessages LoginPressed), Button.attrs [ Spacing.m5 ] ] [ text "Log in" ]
                , loginFormMessageView model
                ]
            ]
        |> Card.view


signupForm model =
    Card.config [ Card.outlinePrimary, Card.attrs [ Spacing.m5 ] ]
        |> Card.headerH1 [] [ text "Registration" ]
        |> Card.block []
            [ Bootstrap.Card.Block.text []
                [ Input.text [ Input.id "usernameInput", Input.onInput (\a -> RegisterMessages (RegUsernameInput a)), Input.value model.registerData.username, Input.attrs [ Spacing.m3, maxWidth100 ] ]
                , Input.password [ Input.id "passwordInput", Input.onInput (\a -> RegisterMessages (RegPasswordInput a)), Input.value model.registerData.password, Input.attrs [ Spacing.m3, maxWidth100 ] ]
                , Button.button [ Button.outlinePrimary, Button.onClick (RegisterMessages RegisterPressed), Button.attrs [ Spacing.m5 ] ] [ text "Register" ]
                , signupFormMessageView model
                ]
            ]
        |> Card.view


loginFormMessageView : Model -> Html Msg
loginFormMessageView model =
    case model.loginData.status of
        InProgress ->
            div [] []

        CredentialsFailure string ->
            Bootstrap.Alert.simpleWarning [ Spacing.m3 ] [ text string ]

        ConnectionFailure ->
            Bootstrap.Alert.simpleWarning [ Spacing.m3 ] [ text "Connection error" ]

        Success ->
            Bootstrap.Alert.simpleSuccess [ Spacing.m3 ] [ text ("Logged in as " ++ model.loginData.username) ]


signupFormMessageView : Model -> Html Msg
signupFormMessageView model =
    case model.registerData.status of
        InProgress ->
            div [] []

        CredentialsFailure string ->
            Bootstrap.Alert.simpleWarning [ Spacing.m3 ] [ text string ]

        ConnectionFailure ->
            Bootstrap.Alert.simpleWarning [ Spacing.m3 ] [ text "Connection error" ]

        Success ->
            Bootstrap.Alert.simpleSuccess [ Spacing.m3 ] [ text ("Signed up as " ++ model.registerData.username) ]


gameView : Model -> Html Msg
gameView model =
    div [] [ gameOverView model, nextPlayerView model, mainGameView model ]


gameOverView : Model -> Html Msg
gameOverView model =
    if model.gameData.result.isOver then
        Grid.row [ Bootstrap.Grid.Row.attrs [ maxWidth100 ] ]
            [ Grid.col []
                [ Button.button [ Button.outlineInfo, Button.disabled True, Button.attrs [ Spacing.m4, Spacing.p4 ] ] [ Html.h4 [] [ text model.gameData.player1, Html.br [] [], text (String.fromInt model.gameData.p1score) ] ]
                ]
            , Grid.col []
                [ Button.button [ Button.outlineInfo, Button.disabled True, Button.attrs [ Spacing.m4, Spacing.p4 ] ] [ Html.h4 [] [ text model.gameData.player2, Html.br [] [], text (String.fromInt model.gameData.p2score) ] ]
                ]
            ]

    else
        div [] []


nextPlayerView : Model -> Html Msg
nextPlayerView model =
    if model.gameData.result.isOver then
        div [] []

    else
        Grid.row []
            [ Grid.col [ Bootstrap.Grid.Col.attrs [ maxWidth100 ] ]
                [ Html.h6 []
                    [ text
                        ("Next is "
                            ++ (if model.gameData.next_player == 1 then
                                    model.gameData.player1

                                else if model.gameData.next_player == 2 then
                                    model.gameData.player2

                                else
                                    " undefined"
                               )
                        )
                    ]
                ]
            ]


mainGameView : Model -> Html Msg
mainGameView model =
    let
        gd =
            model.gameData
    in
    case model.lobbyRole of
        Host ->
            div [ Html.Attributes.id "board", Spacing.m5 ] [ createUpperMancalaRow gd.p2cells gd.p2mancala, createBottomMancalaRow gd.p1cells gd.p1mancala, gameErrorAlert model ]

        Guest ->
            div [ Html.Attributes.id "board", Spacing.m5 ] [ createUpperMancalaRow gd.p1cells gd.p1mancala, createBottomMancalaRow gd.p2cells gd.p2mancala, gameErrorAlert model ]

        None ->
            div [ Html.Attributes.id "board", Spacing.m5 ] []


createUpperMancalaRow : List Int -> Int -> Html Msg
createUpperMancalaRow cells mancala =
    Grid.row []
        [ Grid.col []
            [ ButtonGroup.buttonGroup [ ButtonGroup.large ]
                (List.reverse
                    ([ mancalaCellPlaceholder ] ++ List.map (\score -> ButtonGroup.button [ Button.outlineWarning ] [ text (String.fromInt score) ]) cells ++ [ ButtonGroup.button [ Button.danger ] [ text (String.fromInt mancala) ] ])
                )
            ]
        ]


createBottomMancalaRow : List Int -> Int -> Html Msg
createBottomMancalaRow cells mancala =
    Grid.row []
        [ Grid.col []
            [ ButtonGroup.buttonGroup [ ButtonGroup.large ]
                ([ mancalaCellPlaceholder ] ++ List.map (\indexedScore -> ButtonGroup.button [ Button.outlinePrimary, Button.onClick (GameMsg (MakeTurn (Tuple.first indexedScore))) ] [ text (String.fromInt (Tuple.second indexedScore)) ]) (Array.toIndexedList (Array.fromList cells)) ++ [ ButtonGroup.button [ Button.success ] [ text (String.fromInt mancala) ] ])
            ]
        ]


mancalaCellPlaceholder : ButtonGroup.ButtonItem Msg
mancalaCellPlaceholder =
    ButtonGroup.button [ Button.disabled True, Button.dark, Button.outlineDark, Button.attrs [ Spacing.px4 ] ] [ text "  " ]


gameErrorAlert : Model -> Html Msg
gameErrorAlert model =
    case model.gameData.errorMessage of
        Just string ->
            Bootstrap.Alert.simpleWarning [ Spacing.m5 ] [ text string ]

        Nothing ->
            div [] []

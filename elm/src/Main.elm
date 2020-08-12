module Main exposing (..)

import Bootstrap.Alert
import Bootstrap.Button as Button
import Bootstrap.CDN
import Bootstrap.Card as Card
import Bootstrap.Navbar as Navbar
import Bootstrap.Card.Block
import Bootstrap.Form.Input as Input
import Bootstrap.Grid as Grid
import Bootstrap.Utilities.Spacing
import Browser
import Html exposing (Html, b, div, pre, text)
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

type alias RegisterRequestStatus = LoginRequestStatusModel

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

type alias Username =
    String


type MainWindow
    = LoginPage
    | SignUpPage
    | LobbySelectionPage
    | InGamePage


type alias Model =
    { page : MainWindow
    , loginData : LoginPageModel
    , registerData : RegisterPageModel
    , username : Maybe Username
    , navbarState: Navbar.State
    }


init : () -> ( Model, Cmd Msg )
init _ =
    let
        ( navbarState, navbarCmd ) = Navbar.initialState NavbarMsg
    in
    ( {page=LoginPage, loginData = LoginPageModel "" "" InProgress, registerData = LoginPageModel "" "" InProgress, username = Nothing,navbarState = navbarState}, navbarCmd )

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

type SelectPageMsg
    = SelectLogin
    | SelectRegister
    | SelectLobby
    | SelectInGame

type Msg
    = LoginMessages LoginMessages
    | RegisterMessages RegisterMessages
    | NavbarMsg Navbar.State
    | SelectPageMsg MainWindow
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
                            let loginDataRec=model.loginData in( { model | loginData={loginDataRec |status = Success}}, Cmd.none )

                        CredentialsFailure str ->
                            let loginDataRec=model.loginData in( { model | loginData={loginDataRec |status = CredentialsFailure str }}, Cmd.none )

                        InProgress ->
                            ( model, Cmd.none )

                        ConnectionFailure ->
                            ( model, Cmd.none )

                Err _ ->
                    let loginDataRec=model.loginData in( { model | loginData={loginDataRec |status = ConnectionFailure}}, Cmd.none )

        LoginPressed ->
            ( model
            , Http.post
                { url = "http://localhost:1337/api/user/login"
                , body = Http.stringBody "application/x-www-form-urlencoded" ("username=" ++ model.loginData.username ++ "&password=" ++ model.loginData.password)
                , expect = Http.expectString (\a -> LoginMessages (GotLoginResp a))
                }
            )

        UsernameInput string ->
            let loginDataRec=model.loginData in( { model | loginData={loginDataRec |username = string}}, Cmd.none )


        PasswordInput string ->
            let loginDataRec=model.loginData in( { model | loginData={loginDataRec |password=string}}, Cmd.none )

updateRegister : RegisterMessages -> Model -> ( Model, Cmd Msg )
updateRegister msg model =
    case msg of
        GotRegisterResp result ->
            case result of
                Ok message ->
                    case decodeAuthResponse message of
                        Success ->
                            let registerDataRec=model.registerData in ( { model | registerData={registerDataRec |status = Success}}, Cmd.none )

                        CredentialsFailure str ->
                            let registerDataRec=model.registerData in ( { model | registerData={ registerDataRec |status = CredentialsFailure str }}, Cmd.none )

                        InProgress ->
                            ( model, Cmd.none )

                        ConnectionFailure ->
                            ( model, Cmd.none )

                Err _ ->
                    let registerDataRec=model.registerData in ( { model | registerData={registerDataRec |status = ConnectionFailure}}, Cmd.none )

        RegisterPressed ->
            ( model
            , Http.post
                { url = "http://localhost:1337/api/user/new"
                , body = Http.stringBody "application/x-www-form-urlencoded" ("username=" ++ model.registerData.username ++ "&password=" ++ model.registerData.password)
                , expect = Http.expectString (\a -> RegisterMessages (GotRegisterResp a))
                }
            )

        RegUsernameInput string ->
            let registerDataRec=model.registerData in ( { model | registerData={registerDataRec |username = string}}, Cmd.none )


        RegPasswordInput string ->
            let registerDataRec=model.registerData in ( { model | registerData={registerDataRec |password=string}}, Cmd.none )



update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        LoginMessages lMsg ->
            updateLogin lMsg model


        RegisterMessages rMsg ->
            updateRegister rMsg model

        NavbarMsg state->
            ({model | navbarState = state},Cmd.none)

        SelectPageMsg spm->
            ({model|page=spm},Cmd.none)

        LogOut->
            init ()




-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Navbar.subscriptions model.navbarState NavbarMsg



-- VIEW


view : Model -> Html Msg
view model =
    preprocessView model


preprocessView : Model -> Html Msg
preprocessView model =
    Grid.container [] [ Bootstrap.CDN.stylesheet, navBar model,processView model ]

navBar : Model -> Html Msg
navBar model =
    case model.username of
        Nothing ->
            Navbar.config NavbarMsg
                |> Navbar.withAnimation
                |> Navbar.brand [] [text "Mancala"]
                |> Navbar.items
                    [ Navbar.itemLink [] [Button.button [Button.onClick (SelectPageMsg LoginPage)][text "Log in"]]
                    , Navbar.itemLink [] [Button.button [Button.onClick (SelectPageMsg SignUpPage)][text "Register"]]
                    ]
                |> Navbar.view model.navbarState

        Just username ->
            Navbar.config NavbarMsg
                |> Navbar.withAnimation
                |> Navbar.secondary
                |> Navbar.brand [] [text "Mancala"]
                |> Navbar.items
                   [ Navbar.itemLink [] [text username]
                   , Navbar.itemLink [] [Button.button [Button.onClick LogOut][text "Log out"]]
                   ]
                |> Navbar.view model.navbarState



processView : Model -> Html Msg
processView model =
    case model.page of
        LoginPage ->
            loginForm model

        SignUpPage ->
            signupForm model


        LobbySelectionPage ->
            text "Lobby selection placeholder"


        InGamePage ->
            text "Game placeholder"


loginForm : Model -> Html Msg
loginForm model =
    Card.config [ Card.outlinePrimary, Card.attrs [ Bootstrap.Utilities.Spacing.m5 ] ]
        |> Card.headerH1 [] [ text "Log in to proceed" ]
        |> Card.block []
            [ Bootstrap.Card.Block.text []
                [ Input.text [ Input.id "usernameInput", Input.onInput (\a -> LoginMessages (UsernameInput a)), Input.value model.loginData.username, Input.attrs [ Bootstrap.Utilities.Spacing.m3 ] ]
                , Input.password [ Input.id "passwordInput", Input.onInput (\a -> LoginMessages (PasswordInput a)), Input.value model.loginData.password, Input.attrs [ Bootstrap.Utilities.Spacing.m3 ] ]
                , Button.button [ Button.outlinePrimary, Button.onClick (LoginMessages LoginPressed), Button.attrs [ Bootstrap.Utilities.Spacing.m5 ] ] [ text "Log in" ]
                , loginFormMessageView model
                ]
            ]
        |> Card.view


signupForm model =
    Card.config [ Card.outlinePrimary, Card.attrs [ Bootstrap.Utilities.Spacing.m5 ] ]
        |> Card.headerH1 [] [ text "Registration" ]
        |> Card.block []
            [ Bootstrap.Card.Block.text []
                [ Input.text [ Input.id "usernameInput", Input.onInput (\a -> RegisterMessages (RegUsernameInput a)), Input.value model.registerData.username, Input.attrs [ Bootstrap.Utilities.Spacing.m3 ] ]
                , Input.password [ Input.id "passwordInput", Input.onInput (\a -> RegisterMessages (RegPasswordInput a)), Input.value model.registerData.password, Input.attrs [ Bootstrap.Utilities.Spacing.m3 ] ]
                , Button.button [ Button.outlinePrimary, Button.onClick (RegisterMessages RegisterPressed), Button.attrs [ Bootstrap.Utilities.Spacing.m5 ] ] [ text "Register" ]
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
            Bootstrap.Alert.simpleWarning [ Bootstrap.Utilities.Spacing.m3 ] [ text string ]

        ConnectionFailure ->
            Bootstrap.Alert.simpleWarning [ Bootstrap.Utilities.Spacing.m3 ] [ text "Connection error" ]

        Success ->
            Bootstrap.Alert.simpleSuccess [ Bootstrap.Utilities.Spacing.m3 ] [ text ("Logged in as " ++ model.loginData.username) ]
signupFormMessageView : Model -> Html Msg
signupFormMessageView model =
    case model.registerData.status of
        InProgress ->
            div [] []

        CredentialsFailure string ->
            Bootstrap.Alert.simpleWarning [ Bootstrap.Utilities.Spacing.m3 ] [ text string ]

        ConnectionFailure ->
            Bootstrap.Alert.simpleWarning [ Bootstrap.Utilities.Spacing.m3 ] [ text "Connection error" ]

        Success ->
            Bootstrap.Alert.simpleSuccess [ Bootstrap.Utilities.Spacing.m3 ] [ text ("Signed up as " ++ model.registerData.username) ]

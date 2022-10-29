module Pre exposing (..)

main : Program () () ()
main = 
    Platform.worker 
        { init = \_ -> ((), Cmd.none)
        , update = \_ _ -> ((), Cmd.none)
        , subscriptions = \_ -> Sub.none
        }

port module Main exposing (run)

import M exposing(run)

port send : String -> Cmd msg

main : Program () () ()
main = 
    Platform.worker 
        { init = \_ -> ((), send (M.run holeArgs))
        , update = \_ _ -> ((), Cmd.none)
        , subscriptions = \_ -> Sub.none
        }

holeArgs : List String
holeArgs = []

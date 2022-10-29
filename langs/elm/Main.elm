port module Main exposing (main)

import M

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

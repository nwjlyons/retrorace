module Robots exposing (robot)

import Html exposing (Html)
import Html.Attributes exposing (style)
import Svg exposing (..)
import Svg.Attributes exposing (..)
import Color


white =
    "#d3d3d3"


robot : String -> Float -> Float -> String -> Html msg
robot colour toptop leftleft name =
    let
        eye xx yy =
            rect [ x xx, y yy, width "6", height "6", fill "#2b2b2b" ] []

        pupil xx yy =
            rect [ x xx, y yy, width "2", height "2", fill white ] []

        tooth xx yy =
            rect [ x xx, y yy, width "2", height "5", fill white ] []
    in
        Svg.svg
            [ Html.Attributes.style
                [ ( "bottom", (toString toptop) ++ "vh" )
                , ( "left", (toString leftleft) ++ "%" )
                ]
            , width "10vh"
            , height "10vh"
            , viewBox "0 0 16 16"
            , class "robot"
            ]
            [ rect [ width "24", height "24", fill colour ] []
            , eye "1" "3"
            , eye "9" "3"
            , pupil "3" "5"
            , pupil "11" "5"
            , text_ [ fontFamily "monospace", textAnchor "middle", fontSize "5px", fill white, y "15", x "8" ] [ text name ]
            ]

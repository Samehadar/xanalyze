﻿namespace Route4MeSdk.FSharp

open System
open System.Collections.Generic
open FSharpExt
open FSharp.Data
open Newtonsoft.Json
open Newtonsoft.Json.Linq

[<CLIMutable>]
type AvoidanceZone = {
    [<JsonProperty("territory_id")>]
    Id : string

    [<JsonProperty("member_id")>]
    MemberId : int option

    [<JsonProperty("territory_name")>]
    Name : string

    [<JsonProperty("territory_color")>]
    Color : string

    [<JsonProperty("addresses")>]
    Addresses : string[]
    
    [<JsonProperty("territory")>]
    Parameters : TerritoryParameters }

    with
        static member Get(?take:int, ?skip:int, ?apiKey) =
            let query = 
                [ take |> Option.map(fun v -> "limit", v.ToString())
                  skip |> Option.map(fun v -> "offset", v.ToString()) ]
                |> List.choose id

            Api.Get(Url.V4.avoidance, [], query, apiKey)
            |> Result.map Api.Deserialize<AvoidanceZone[]>

        static member Get(territoryId, ?apiKey) =
            let query = [("territory_id", territoryId)]

            Api.Get(Url.V4.avoidance, [], query, apiKey)
            |> Result.map Api.Deserialize<AvoidanceZone>

        static member Add(name:string, color:string, parameters:TerritoryParameters, ?apiKey) =
            let request = 
                [("territory_name", box name)
                 ("territory_color", box color)
                 ("territory", box parameters)]
                |> dict

            Api.Post(Url.V4.avoidance, [], [], apiKey, request)
            |> Result.map Api.Deserialize<AvoidanceZone>

        static member Update(zone:AvoidanceZone, ?apiKey) =
            zone.Id
            |> Option.ofObj
            |> Result.ofOption(ValidationError("AvoidanceZone Id must be supplied."))
            |> Result.andThen(fun id ->
                Api.Post(Url.V4.avoidance, [], [], apiKey, zone)
                |> Result.map Api.Deserialize<AvoidanceZone>)

        member self.Update(?apiKey) =
            match apiKey with
            | None -> AvoidanceZone.Update(self)
            | Some v -> AvoidanceZone.Update(self, v)

        static member Delete(zone:AvoidanceZone, ?apiKey) =
            zone.Id
            |> Option.ofObj
            |> Result.ofOption(ValidationError("AvoidanceZone Id must be supplied."))
            |> Result.andThen(fun id ->
                let query = [("territory_id", id)]
                Api.Delete(Url.V4.avoidance, [], query, apiKey))

        member self.Delete(?apiKey) =
            match apiKey with
            | None -> AvoidanceZone.Delete(self)
            | Some v -> AvoidanceZone.Delete(self, v)

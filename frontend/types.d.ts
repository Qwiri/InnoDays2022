interface Game {
    ID: number
    StartTime: string
    EndTime: {
        Time: string
        Valid: boolean
    }
    ScoreBlack: number
    ScoreWhite: number
    KickaeID: number
    Kickae: Kickae
    UpdatedAt: string
    Players: Array<Player>
    Goals: Array<Goal>
}

interface Player {
    PlayerID: string
    Player: {
        ID: string,
        Nick: string,
        Elo: number
    }
    Team: number

}
interface Kickae {
    ID: number
    Room: string
    Note: string
}
interface Goal {
    ID: number,
    Team: number,
    Time: string
}
interface Pending {
    Player: {
        ID: string,
        Nick: string,
        Elo: number
    }
    Pending: {
        Team: number
        AddedAt: string
    }
}
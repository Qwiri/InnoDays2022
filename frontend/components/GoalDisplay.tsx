import { Fragment } from "react"
import styles from "../styles/goalDisplay.module.scss"

export default function GoalDisplay({goals, startTime}: {goals: Array<Goal>, startTime: string}) {

    let startTimeStamp = Date.parse(startTime)


    // return goaltime as string: `mm ss`
    function formatTime(time: string) {

        let goalTime = Date.parse(time)
        let diff = goalTime - startTimeStamp

        let minutes = Math.floor(diff / 60000)
        let minutesString = minutes.toString().length < 2 ? "0"+minutes.toString() : minutes.toString()
        let seconds = Math.floor((diff / 1000) % 60)
        let secondsString = seconds.toString().length < 2 ? "0"+seconds.toString() : seconds.toString()

        return `${minutesString}m ${secondsString}s`

    }

    function paintCircle(g: Goal, team: number) {
        let circleColor = (g.Team == team ? 'bg-white' : 'bg-stone-500');
        let offset = (team == 2 ? 'after:top-1/2' : 'after:bottom-1/2')

        return(<span className={`${styles.preserve3d} ${circleColor} relative w-4 h-4 rounded-full bg-white text-center after:absolute after:translate-z-1 after:w-1 after:h-full after:left-[calc(50%-0.1rem)] after:bg-stone-600 ${offset}`}></span>)
    }

    return(
        <div className="grid grid-flow-row grid-cols-3 content-start justify-items-center items-center gap-2">

            {goals.map((g, i) => {
                return(
                    <Fragment key={i}>
                        {paintCircle(g, 1)}
                        <p className='h-4 leading-none'>{formatTime(g.Time)}</p>
                        {paintCircle(g, 2)}
                    </Fragment>
                )
            })}
        </div>
    )

}
import { Fragment } from "react";
import styles from "../styles/goalDisplay.module.scss";

export default function GoalDisplay({goals, startTime}: {goals: Array<Goal>, startTime: string}) {

	const startTimeStamp = Date.parse(startTime);



	function paintCircle(g: Goal, team: number) {
		const circleColor = (g.Team == team ? "bg-white" : "bg-stone-500");
		const offset = (team == 2 ? "after:top-1/2" : "after:bottom-1/2");

		return(<span className={`${styles.preserve3d} ${circleColor} relative w-4 h-4 rounded-full bg-white text-center after:absolute after:translate-z-1 after:w-1 after:h-full after:left-[calc(50%-0.1rem)] after:bg-stone-600 ${offset}`}></span>);
	}

	return(
		<div className="grid grid-flow-row grid-cols-3 content-start justify-items-center items-center gap-2">

			{goals.map((g, i) => {
				return(
					<Fragment key={i}>
						{paintCircle(g, 1)}
						<p className='h-4 leading-none'>{formatTime(startTimeStamp, g.Time)}</p>
						{paintCircle(g, 2)}
					</Fragment>
				);
			})}
		</div>
	);

}
// return goaltime as string: `mm ss`
export function formatTime(startTimeStamp: number, time: string) {

	const goalTime = Date.parse(time);
	const diff = goalTime - startTimeStamp;

	const minutes = Math.floor(diff / 60000);
	const minutesString = minutes.toString().length < 2 ? "0"+minutes.toString() : minutes.toString();
	const seconds = Math.floor((diff / 1000) % 60);
	const secondsString = seconds.toString().length < 2 ? "0"+seconds.toString() : seconds.toString();

	return `${minutesString}m ${secondsString}s`;

}
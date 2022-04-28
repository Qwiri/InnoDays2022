import styles from "/styles/kickae.module.scss";
import { useEffect, useState } from "react";
import TeamCard from "../../components/TeamCard";
import GoalDisplay, { formatTime } from "../../components/GoalDisplay";

import { QRCodeSVG } from "qrcode.react";

export default function KickaePage({ id }: {id: number}) {

	const [scoreWhite, setScoreWhite] = useState("?");
	const [scoreBlack, setScoreBlack] = useState("?");
	const [teamWhite, setTeamWhite] = useState([]);
	const [teamBlack, setTeamBlack] = useState([]);
	const [goals, setGoals] = useState(Array<Goal>());
	const [game, setGame] = useState();

	const [qrCode, setQrCode] = useState("");
	const [gameTime, setGameTime] = useState("");

	async function fetchData() {
		const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/p/monitor/${id}`).catch(() => {return;});
		const json = await res.json();
        
		const teamWhite: Array<Player> = [];
		const teamBlack: Array<Player> = [];
		let game: Game;

		if (json["Game"] === null) {
			if (json["Pending"] !== null) {

				let players: Array<Pending>;
				players = json["Pending"];
				players.forEach((p) => {
					if (p.Pending.Team === 1) {
						teamBlack.push(p);
					} else if (p.Pending.Team === 2) {
						teamWhite.push(p);
					}
				});
			}

		} else {
			game = json.Game;
			setGameTime(formatTime(Date.parse(game.StartTime), new Date(Date.now()).toISOString()));

			const players = game.Players;
			players.forEach((p) => {
				if (p.Team === 1) {
					teamBlack.push(p);
				} else if (p.Team === 2) {
					teamWhite.push(p);
				}
			});
		}

		setTeamWhite(teamWhite);
		setTeamBlack(teamBlack);
		setGame(game);
		setGoals(game?.Goals ?? []);
		setScoreBlack(game?.ScoreBlack?.toString() ?? "?");
		setScoreWhite(game?.ScoreWhite?.toString() ?? "?");

	}

	useEffect(() => {
		setQrCode(window.location.href);
		const interval = setInterval(() => {
			fetchData();
		}, 1000);

		// clear interval on re-render to avoid memory leaks
		return () => clearInterval(interval);
	}, []);

	return(
		<>
			<div className="absolute top-0 left-0 w-full overflow-hidden">
				<svg className="relative block h-16 w-full" data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 120" preserveAspectRatio="none">
					<path d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V0H0V27.35A600.21,600.21,0,0,0,321.39,56.44Z" className="fill-white"></path>
				</svg>
			</div>
			<div className="p-8 mt-16 flex gap-4 text-4xl">
				{/* QR-code */}
				<QRCodeSVG size={64} value={qrCode} />
				<div className='flex-column'>
					<h1>Game in progress</h1>
					{game?.ID === undefined ? null : <h3 className='text-gray-500 text-2xl'>#{game.ID}</h3> }
				</div>
			</div>
			{/* wrapper */}
			<div className='flex justify-around'>
				{/* Team white */}
				<div className='flex flex-col gap-8'>

					{/* Point display */}
					<div className='w-64 h-64 rounded-lg flex justify-center items-center bg-gradient-to-b from-gray-200 via-gray-400 to-gray-600'>
						<h1 className='text-8xl font-bold'>{scoreWhite}</h1>
					</div>

					{/* user cards */}
					<TeamCard team={teamWhite} />

				</div>

				{/* goal display */}
				{(goals.length != 0) ? <GoalDisplay goals={goals} startTime={game.StartTime} /> : null}

				{/* Black team */}
				<div className='flex flex-col gap-8'>

					{/* Point display */}
					<div className='w-64 h-64 rounded-lg flex justify-center items-center bg-gradient-to-b from-gray-200 via-gray-400 to-gray-600'>
						<h1 className='text-8xl font-bold text-black'>{scoreBlack}</h1>
					</div>

					{/* user cards */}
					<TeamCard team={teamBlack} />

				</div>

			</div>

			{/* current time */}
			<h1 className="text-4xl text-center">{gameTime}</h1>

                
		</>
	);
}
export async function getServerSideProps(ctx) {
	const { id } = ctx.query;
	return {
		props: {
			id,
		},
	};
}
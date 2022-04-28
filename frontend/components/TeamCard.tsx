import Image from "next/image";

export default function TeamCard({team}: {team: Array<Player>}) {
	return(
		<>
			{team.map((p, i) => {
				return(
					<div key={i} className='flex-col'>
						<div className='bg-zinc-800 flex p-2 rounded-md gap-2'>
							<Image height="48" width="48" src={`https://avatars.dicebear.com/api/bottts/${p.Player.ID}.svg`} />
							<div className='flex-col'>
								<h1 className='font-bold' >{p.Player.Nick !== "" ? p.Player.Nick : p.Player.ID }</h1>
								<h3 className='text-slate-600'>{`${p.Player.Elo} Elo`}</h3>
							</div>
						</div>
					</div>
				);
			})}
		</>
	);
}
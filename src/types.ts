export enum Options {
	SearchMal,
	GetList,
	SearchList,
	Add,
	Remove,
	ModifyCompletion,
	ModifyStart,
	ImportMal,
	ImportKitsu,
	ImportHianime,
	ImportHaru,
	Completion,
	Stats,
}

export enum Completion {
	PlanToWatch,
	OnHold,
	Dropped,
	Watching,
	Completed,
}

export interface EntryData {
	name: string;
	mal_id: number;
	completion: Completion;
	start_date: string;
}

export interface MalAnime {
	series_animedb_id: number;
	series_title: string;
	my_start_date?: string;
	my_status: string;
}

export interface MalFormat {
	myanimelist: {
		anime: MalAnime[];
	};
}

export interface HiAnimeData {
	link: string;
	name: string;
}

export interface HiAnimeFormat {
	[key: string]: HiAnimeData[];
}

interface SearchData {
	mal_id: number;
	title: string;
	title_japanese: string;
	score: number;
	episodes: number;
	aired: {
		string: string;
	};
}

export interface SearchResult {
	data: SearchData[];
}

export interface AnimeResult {
	data: {
		title: string;
	};
}

export interface Config {
	list_location: string;
	colors: boolean;
	database_url?: string;
}

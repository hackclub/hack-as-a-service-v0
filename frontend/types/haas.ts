export interface IUser {
  ID: string;
  Name: string;
  Avatar: string;
  SlackUserID: string;
}

export interface ITeam {
  ID: number;
  Name: string;
  Avatar: string;
  Automatic: boolean;
  Personal: boolean;
  Apps: IApp[];
}

export interface IApp {
  ID: number;
  Name: string;
  ShortName: string;
  TeamID: number;
}

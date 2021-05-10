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
  Expenses: string;
  Users: IUser[];
}

export interface IApp {
  ID: number;
  Name: string;
  ShortName: string;
  TeamID: number;
}

export interface ILetsEncrypt {
  LetsEncryptEnabled: boolean;
}

export interface IBuild {
  ID: string;
  ExecID: string;
  AppID: number;
  StartedAt: number;
  EndedAt: number;
  Running: boolean;
  Events: string[];
  Status: number;
}

export type KVConfig = {
  [id: string]: {
    key: string;
    keyEditable: boolean;
    valueEditable: boolean;
    obscureValue: boolean;
    value: string;
  };
};

export interface IAddon {
  name: string;
  activated: boolean;
  description: string;
  img: string;
  id: string;
  config: KVConfig;
  storage: string;
  price: string;
}

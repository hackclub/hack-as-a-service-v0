import { IAddon } from "../types/haas";

export let devAddons: IAddon[] = [
  {
    price: "3",
    id: "0",
    name: "PostgreSQL",
    activated: false,
    img: "https://upload.wikimedia.org/wikipedia/commons/2/29/Postgresql_elephant.svg",
    description:
      "PostgreSQL is the best database to ever exist. Sometimes it's hard to understand why people use other databases like <insert DB here>. Postgres >>>",
    storage: "1.2 GB",
    config: {
      ddfddfy47: {
        key: "USER",
        value: "ROOT",
        valueEditable: true,
        keyEditable: true,
        obscureValue: false,
      },
      djfusdhf8e74: {
        key: "PASSWORD",
        value: "iL0V3dA1a",
        valueEditable: true,
        keyEditable: false,
        obscureValue: true,
      },
    },
  },
  {
    price: "3",
    id: "1",
    name: "MongoDB",
    activated: true,
    storage: "3.6 GB",
    img: "https://media-exp1.licdn.com/dms/image/C560BAQGC029P7UbAMQ/company-logo_200_200/0/1562088387077?e=2159024400&v=beta&t=lEY4Obku1xJ3BB_BpN3Np9ILy8_zaB1_yjsfH9A57qs",
    description:
      "MongoDB is the best database to ever exist. Sometimes it's hard to understand why people use other databases like <insert DB here>. MongoDB >>>",
    config: {
      u488h: {
        key: "ADMIN_USER",
        value: "root",
        keyEditable: false,
        valueEditable: true,
        obscureValue: false,
      },
      hfofs9: {
        key: "PASSWORD",
        value: "uwuowo123",
        keyEditable: false,
        valueEditable: true,
        obscureValue: true,
      },
    },
  },
  {
    price: "3",
    id: "2",
    name: "Redis",
    activated: false,
    storage: "1.9 GB",
    img: "https://www.nditech.org/sites/default/files/styles/small_photo/public/redis-logo.png?itok=LrULOkWT",
    description:
      "Redis is the best database to ever exist. Sometimes it's hard to understand why people use other databases like <insert DB here>. Redis >>>",
    config: {
      ddfsfyy5: {
        key: "ADMIN_USER",
        value: "root",
        keyEditable: false,
        valueEditable: true,
        obscureValue: false,
      },
      idsf8: {
        key: "PASSWORD",
        value: "uwuowo123",
        keyEditable: false,
        valueEditable: true,
        obscureValue: true,
      },
    },
  },
];

export const devAddonsOriginal: IAddon[] = devAddons;

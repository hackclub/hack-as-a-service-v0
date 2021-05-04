import AppLayout from "../../../../layouts/app";
import {useRouter} from "next/router"
import {devAddons, IAddon} from "./index"
export default function AppAddonDetails() {
  const router = useRouter();
  const { addon } = router.query;
  const a = devAddons.find(element => element.id == addon);
  
return <AppLayout selected="Dashboard">{a ? <Content {...a}/> : "lol wyd chief dis addon dont exist"}</AppLayout>;
}

// TODO: Make component good lol
function Content({id, name, activated, description}:IAddon) {
return <p>{name} - {description}</p>
}
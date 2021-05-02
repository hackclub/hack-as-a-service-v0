import { useEffect } from "react";
import AppLayout from "../../../layouts/app";

export default function AppDashboardPage() {
  useEffect(() => {
    console.log("mount");

    return () => {
      console.log("unmount");
    };
  }, []);

  return <AppLayout>app logs</AppLayout>;
}

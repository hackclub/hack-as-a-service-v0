import DashboardLayout from "../../layouts/dashboard";
import { useRouter } from "next/router";

export default function TeamPage() {
  const router = useRouter();
  const { id } = router.query;

  return (
    <DashboardLayout title={`Team ${id}`} sidebarSections={[]}>
      hi there
    </DashboardLayout>
  );
}

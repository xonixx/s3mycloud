import { createRouter, createWebHistory } from "vue-router";
// import HomeView from "../views/HomeView.vue";
import ListFilesView from "../views/ListFilesView.vue";
import FilesUploadView from "../components/files-dnd-uploader/FilesUploadView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "list",
      component: ListFilesView,
      props: (route) => ({
        pageSize: asInt(route.query.pageSize),
        page: asInt(route.query.page),
        q: route.query.q,
      }),
    },
    {
      path: "/upload",
      name: "upload",
      component: FilesUploadView,
    },
    {
      path: "/about",
      name: "about",
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("../views/AboutView.vue"),
    },
  ],
});

function asInt(v: any, defaultVal = undefined): number | undefined {
  return parseInt(v) || defaultVal;
}

export default router;

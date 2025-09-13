import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';

import LoginView from '../views/LoginView.vue';
import DashboardView from '../views/DashboardView.vue';
import RegisterView from '../views/RegisterView.vue';
import AdminDashboardView from '../views/AdminDashboardView.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView,
  },
  {
    path : '/admin',
    name : 'AdminDashboard',
    component : AdminDashboardView,
  },
  {
    path: '/', 
    redirect: '/login',
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;

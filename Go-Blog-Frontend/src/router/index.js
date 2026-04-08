import { createRouter, createWebHistory } from 'vue-router'
import HomeContainer from '../views/Home-Container.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: HomeContainer,
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/Login-Page.vue'),
    },
    {
      path: '/article/:id',
      name: 'article-detail',
      component: () => import('../views/Article-Detail.vue'),
    },
    // admin
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/Admin-Layout.vue'),
      redirect: '/admin/articles',
      meta: { requiresAuth: true },
      children: [
        {
          path: 'articles',
          name: 'admin-articles',
          component: () => import('../views/admin/Article-List.vue'),
        },
        {
          path: 'categories',
          name: 'admin-categories',
          component: () => import('../views/admin/Category-Tag.vue'),
        },
        {
          path: 'comments',
          name: 'admin-comments',
          component: () => import('../views/admin/CommentList.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router

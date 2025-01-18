import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ConversationsView from '../views/ConversationsView.vue'
import ChatView from '../views/ChatView.vue'
import NewConversationView from '../views/NewConversationView.vue'
import NewGroupView from '../views/NewGroupView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/conversations'
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/conversations',
      name: 'conversations',
      component: ConversationsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/conversations/:id',
      name: 'chat',
      component: ChatView,
      props: true,
      meta: { requiresAuth: true }
    },
    {
      path: '/new-conversation',
      name: 'newConversation',
      component: NewConversationView,
      meta: { requiresAuth: true }
    },
    {
      path: '/new-group',
      name: 'newGroup',
      component: NewGroupView,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !sessionStorage.getItem('authToken')) {
    next('/login')
  } else {
    next()
  }
})

export default router
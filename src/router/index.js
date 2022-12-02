import Vue from 'vue'
import Router from 'vue-router'
import SM2 from '@/components/SM2'
import SM3 from '@/components/SM3'
import SM4 from '@/components/SM4'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      redirect:'/sm3'
    },
    {
      path: '/sm2',
      name: 'SM2',
      component: SM2
    },
    {
      path: '/sm3',
      name: 'SM3',
      component: SM3
    },
    {
      path: '/sm4',
      name: 'SM4',
      component: SM4
    }
  ]
})

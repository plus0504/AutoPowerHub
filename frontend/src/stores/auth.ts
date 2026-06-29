import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') ?? '')
  const username = ref(localStorage.getItem('username') ?? '')

  function setAuth(newToken: string, user: string) {
    token.value = newToken
    username.value = user
    localStorage.setItem('token', newToken)
    localStorage.setItem('username', user)
  }

  function logout() {
    token.value = ''
    username.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('username')
  }

  return { token, username, setAuth, logout }
})

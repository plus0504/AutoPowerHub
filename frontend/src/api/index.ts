import axios from 'axios'

export interface Device {
  id: number
  name: string
  mac: string
  service_uuid: string
  characteristic_uuid: string
  enabled: boolean
}

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use((cfg) => {
  const token = localStorage.getItem('token')
  if (token) cfg.headers.Authorization = `Bearer ${token}`
  return cfg
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  },
)

export const login = (username: string, password: string) =>
  api.post<{ token: string }>('/login', { username, password })

export const getDevices = () =>
  api.get<{ devices: Device[] }>('/device')

export const powerDevice = (id: number) =>
  api.post<{ message: string }>(`/device/${id}/power`)

export const testDevice = (id: number) =>
  api.post<{ message: string }>(`/device/${id}/test`)

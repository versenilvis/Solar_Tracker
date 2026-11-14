import { browser } from '$app/environment'

export type SensorData = {
  volt: number
  curr: number
  power: number
  energy: number
  azi: number
  ele: number
  ldr: [number, number, number, number]
  mode: string
  timestamp: string
}

export type DailyStats = {
  max_power: number
  total_energy: number
}

type WsState = {
  connected: boolean
  latest: SensorData | null
  history: SensorData[]
  stats: DailyStats
}

const MAX_HISTORY = 60

function createWsStore() {
  let state = $state<WsState>({
    connected: false,
    latest: null,
    history: [],
    stats: {
      max_power: 0,
      total_energy: 0
    }
  })

  let socket: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null

  async function fetchStats() {
    try {
      const res = await fetch('/api/stats')
      const data = await res.json()
      state.stats = data
    } catch (e) {
      console.error('failed to fetch stats', e)
    }
  }

  function connect() {
    if (!browser) return

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    socket = new WebSocket(`${protocol}//${window.location.host}/ws`)

    socket.onopen = () => {
      state.connected = true
      fetchStats()
      if (reconnectTimer) {
        clearTimeout(reconnectTimer)
        reconnectTimer = null
      }
    }

    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      if (msg.type === 'history') {
        state.history = msg.data.slice(-MAX_HISTORY)
        if (state.history.length > 0) {
          state.latest = state.history[state.history.length - 1]
        }
      } else if (msg.type === 'live') {
        const data: SensorData = msg.data
        state.latest = data
        state.history = [...state.history.slice(-(MAX_HISTORY - 1)), data]

        if (data.power > state.stats.max_power) state.stats.max_power = data.power
      }
    }

    socket.onclose = () => {
      state.connected = false
      reconnectTimer = setTimeout(connect, 2000)
    }

    socket.onerror = () => {
      socket?.close()
    }
  }

  connect()

  return {
    get connected() { return state.connected },
    get latest() { return state.latest },
    get history() { return state.history },
    get stats() { return state.stats },
  }
}

export const ws = createWsStore()

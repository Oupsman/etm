export interface UserSession {
  authorized: boolean,
  exp: number,
  iss: string,
  role: string,
  sub: number,
  uuid: string,
  token: string,
}

export interface User {
  ID: number,
  username: string,
  active_category_id: number,
  isadmin: string,
  oidc_subject?: string,
  oidc_provider?: string,
}

export interface Device {
  ID: number
  UUID: string
  CreatedAt: string
  userid: number
  devicename: string
  api_key_prefix: string
  trusted: boolean
  lastusedat: string
}

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
}

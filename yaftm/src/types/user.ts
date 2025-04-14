export interface UserSession {
  authorized: boolean,
  exp: number,
  iss: string,
  role: string,
  sub: number,
  uuid: string,
  token: string,
}

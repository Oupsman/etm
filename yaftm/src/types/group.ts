export interface Group {
  ID: number
  name: string
  owner_id: number
  CreatedAt: string
}

export interface GroupMember {
  user_id: number
  username: string
  email: string
  role: string
}

export interface GroupDetail {
  group: Group
  members: GroupMember[]
}

export interface UserSearchResult {
  ID: number
  username: string
  email: string
  isadmin: string
}

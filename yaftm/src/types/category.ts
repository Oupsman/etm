export interface Category {
  ID: number
  name: string
  color: string
  userid?: number
  active?: boolean
  shared?: boolean
  group_id?: number
  group_name?: string
  user_role?: string  // 'reader' | 'writer' | 'owner' — empty when the user owns the category
}

export interface NewCategory {
  name: string
  color: string
}

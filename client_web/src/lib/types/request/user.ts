export interface CreateUserRequest {
	username: string;
	password: string;
	hi_res: boolean;
}

export interface UpdateUserRequest {
	username?: string;
	password?: string;
	hi_res?: boolean;
}

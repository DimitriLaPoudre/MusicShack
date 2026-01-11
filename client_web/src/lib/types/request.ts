export interface RequestUserLogin {
	username: string;
	password: string;
}

export interface RequestUser {
	username: string;
	password: string;
	hiRes: boolean;
}

export interface RequestAdmin {
	password: string;
}

export interface RequestAdminPassword {
	oldPassword: string;
	newPassword: string;
}

export interface RequestInstance {
	url: string;
}

export interface RequestFollow {
	provider: string;
	id: string;
}

export interface RequestDownload {
	provider: string;
	type: string;
	id: string;
}

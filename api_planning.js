POST /v1/auth/login  // логин
POST /v1/auth/register // Регистрация
POST /v1/auth/restore // Восстановление аккаунта

GET /v1/users/me // Получение своей информации
PUT /v1/users/me // Редактирование своей информации
GET /v1/users/me/posts // Получение своих постов
GET /v1/users/{userId}/posts // Получение постов пользователя
GET /v1/users/my/comments // Получение постов пользователя
GET /v1/users/{userId}/comments // Получение постов пользователя


GET /v1/posts // Все посты
POST /v1/posts // Создать пост
PUT /v1/posts/{postId} // Изменить пост
DELETE /v1/posts/{postId} //Удалить пост

GET /v1/posts/{postId}/comments //Получить комментарии пользователя к посту
POST /v1/posts/{postId}/comments  // Создать комментарий
PUT /v1/posts/{postId}/comments/{commentId}  // Изменить комментарий
DELETE /v1/posts/{postId}/comments/{commentId}  // Удалить комментарий



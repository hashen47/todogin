-- +goose Up
-- +goose StatementBegin
CREATE TABLE `todos` (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title VARCHAR(100) NOT NULL,
    content VARCHAR(255) NOT NULL,
    user_id INT UNSIGNED,
    done TINYINT(1) DEFAULT 0 NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `todos`; 
-- +goose StatementEnd

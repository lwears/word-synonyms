# 📝 Word Synonyms App

Welcome to the Word Synonyms App! This application allows you to manage words and their synonyms through a simple API. It includes a backend written in Go and a frontend built with TypeScript and React.

## 📂 Project Structure

```md
word-synonyms/
├── backend/ # Go backend code
└── frontend/ # React frontend code
```

## 🚀 Getting Started

### Prerequisites

- [Go 1.22.5](https://golang.org/dl/) (for the backend)
- [Node.js and npm](https://nodejs.org/en/download/) (for the frontend)

### Run API

```sh
cd backend
go run main.go
```

### Run Frontend

```bash
cd frontend
pnpm install
pnpm run dev

#OR
npm install
npm run dev
```

Browse to http://localhost:5173

## :clipboard: Todos

- Backend: Database Transactions
- Frontend: Move URL into env variable
- Eslint & Prettier Github workflow
- [OpenApi TypeScript](https://github.com/openapi-ts/openapi-typescript) (OpenApi TypeScript & Fetch)
- [Validator](https://github.com/go-playground/validator) Validator in struct tags - Go Playground Validator

## Dependencies

### Backend

- [rs/cors](https://github.com/rs/cors) (I chose this for simplicity, but it would have been nice to implement it myself)

### Frontend

- [KY HttpClient](https://github.com/sindresorhus/ky) (I have wanted to try this for a while.)
- [Sonner](https://sonner.emilkowal.ski/) (Toasts)
- [Zod](https://zod.dev/) (Schema Validation)
- [React Hook form](https://www.react-hook-form.com/) (Form Library)
- [Tailwind CSS](https://tailwindcss.com/)(CSS utility classes)

## 📚 Resources

Here are some useful resources to help you get started with this project:

- [Gophers Discord](https://discord.com/invite/golang)
- [Practical Go](https://practicalgobook.net/posts/go-sqlite-no-cgo/) (for sql integration)
- [API with stdlib 1.22+](https://medium.com/@matteopampana/write-the-perfect-rest-api-with-go-1-22-fc7d510230c4) (Build API with standard library)
- [Github Actions](https://olegk.dev/github-actions-and-go#heading-github-actions)
- [Alex Edwards Blog](https://www.alexedwards.net/blog/organising-database-access) (Organising DB Access)
- [Alex Edwards Blog](https://www.alexedwards.net/blog/introduction-to-using-sql-databases-in-go) (SQL Databases in Go)
- [Alex Edwards Blog](https://www.alexedwards.net/blog/making-and-using-middleware) (EnforceJsonHandler Middleware)
- [React Hook Form](https://www.freecodecamp.org/news/react-form-validation-zod-react-hook-form/) (React-Hook-Form Validation with Zod)
- [Resource 4](#)

Feel free to explore these links to gain a better understanding of the project and its dependencies. Happy coding!

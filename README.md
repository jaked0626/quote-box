# Quote Box  
A place for me to keep my favorite quotes.  

Feature ideas:  
- Topic analysis of quotes per user  
- Network visualization of similarities between user  
- Word Map of posted quotes  

# Hosting Postgres  
- https://console.neon.tech/app/projects/curly-sea-40334847

# TIPS
- All Fatal and Panic error logs should only be called from within main. Other errors should be passed up the stack.
- If all handlers are in the same package, wrap all application dependencies in an application struct (data store) and define handlers as methods to the application.
- If handlers exist in separate packages, use a closure model
- In postgres, all users are always allowed to create tables (and delete their own tables). We can restrict user permissions on specific tables that are not owned by them.
- https://stackoverflow.com/questions/19309416/grant-permissions-to-user-for-any-new-tables-created-in-postgresql
- `curl -i -X POST http://localhost:4000/snippet/create`
- check tables in database using `\dt+`
- `curl -iL -X POST 'http://localhost:4000/snippet/create?title=test&content=test&expires=7'`
- wgo run ./cmd/web -addr=":4000" -dsn="postgresql://web:secret@localhost:5432/snippet_db?sslmode=disable"

- use form-encoding on thunderclient


# TODO: 
- default values for author, work are not working 
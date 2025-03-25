# Order service

Develop the following service using a language of your choice between [Java, PHP, Python].
You're free to use any framework and library, but be prepared to have to explain your choices!

Make the service easy to run either using a script, a Dockerfile or writing the commands in a documentation file.
And make sure to write some tests!

## Context

The team needs a system to monitor and manage daily user orders. Your task is to design and
implement the backend according to the specified requirements.

### Assumptions

Assume there is a frontend that interacts with your APIs and implements the following
functionalities:

1. **Order Viewing Page**: Includes filters for date and the ability to search by name and description.
2. **Detailed Order View**: Displays order information along with associated products.
3. **Order Management**: Allows creating, editing or deleting an order.

### Objectives

#### Stock Management
Implement logic to track and manage product stock levels.
Hint: Think about how stock levels should be managed when orders are created or
modified and how to handle concurrency issues if multiple orders might affect stock
simultaneously. 

#### Efficient Search [Optional]
Enhance order search functionality by integrating tools like
Elasticsearch, Meilisearch or a similar indexing solution to manage efficient searching and
filtering.

## Instructions

We have included a simple compose file to make it easy to run some support services that you could need, but feel free
to change or use something else completely.

To use it you need to have Docker installed, then you can run:

```bash
docker compose up
```

## Miscellaneous

Feel free to add everything that's useful or necessary to write high quality code, such as: automatic tests, static code
analysis, coding standards automations, etc.

## Evaluation

Your solution will be evaluated based on the following criteria:
- Correctness and completeness of the RESTful API implementation
- Code quality and adherence to design principles and best practices in backend development
- Overall solution quality (including documentation, build, tests, ...)

The aim of this test is not only to assess your technical skills but also to give you the opportunity to
express yourself freely. We want to see how you approach a real project, organize your work, and
solve problems.

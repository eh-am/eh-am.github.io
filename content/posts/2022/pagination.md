---
title: "Pagination"
date: 2022-12-30T16:28:02Z
---

I was looking into how to do proper pagination. I already knew about offset based pagination, normally used (among other places)in forums. It works as follows:

* Each page has a fixed number of elements (eg 10), aka `LIMIT` in sql.
* When querying, you use, `OFFSET` to skip the first N elements. To calculate N you use `current page - 1 * number of elements per page`. If current page is 1, 0 * 10 is 0. If current page is 2, we offset 1 * 10 which is 10, skipping the first 10 items.


This approach has the following problems:
* In relational databases, the entire thing needs to be sorted before discarting the unwanted values (source: [use-the-index-luke](https://use-the-index-luke.com/no-offset)). 
* Let's say you are in page 1, then click on page 2, but in the meanwhile a new item was inserted in page 1, what should happen? Well, with this approach the last item in page 1 is shifted to page 2. This behaviour is common in forums forums, eg the post has "fallen to the second page", but when displaying a list of items it may be a bit weird to see the same item twice.

Then I came across the cursor implementation which deals with these problems:
* Each page returns a `cursor`, which is really the sort key of the last value (in our example, if we are sorting by id, last item would be 10)
* Then, when querying the database, we use `WHERE ID > 10`
* Obviously, when querying the next page, you pass that cursor so that the backend knows where to start.

Keep in mind this doesn't really work when you want to know all pages upfront, for example when one may arbitrarily go to page 6.



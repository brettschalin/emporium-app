
/*  PLUMBING  */

async function request(url, method, query = {}, body) {
    try {
        const queryKeys = Object.keys(query);
        if (queryKeys.length) {
            url += "?";
            queryKeys.forEach(function(key, index) {
                if (index != 0) {
                    url += "&";
                }
                const value = query[key];
                url += encodeURIComponent(key) + "=" + encodeURIComponent(value);
            })
        }
    } catch (e) {
        console.error(e);
    }

    const options = {
        method: method || "GET",
    }
    if (body) {
        options.body = JSON.stringify(body)
        options.headers = {
            "Content-Type": "application/json"
        }
    }
    
    let response = await fetch(url, options)
    try {
        return await(response.json());
    } catch(e) {
        console.error(e);
        return {};
    }
}

async function GET(url, query) {
    return await request(url, "GET", query);
}

async function POST(url, query, body) {
    return await request(url, "POST", query, body);
}

async function PUT(url, query, body) {
    return await request(url, "PUT", query, body);
}
async function DELETE(url, query) {
    return await request(url, "DELETE", query);
}


/*    REQUEST METHODS    */

async function getPosts(in_mod_queue) {
    return await GET("/api/posts/", {in_mod_queue});
}

async function createUser(name, email, isAdmin) {
    return await POST("/api/users/", undefined, {
        name,
        email,
        is_admin: isAdmin || false
    })
}

async function createPost(created_by, contents, root_parent, direct_parent) {
    return await POST("/api/posts/", undefined, {
        created_by,
        contents,
        root_parent,
        direct_parent,
    });
}

async function approvePosts(approved, ids) {
    return await POST("/api/posts/approve", undefined, {
        approved: !!approved,
        ids,
    })
}
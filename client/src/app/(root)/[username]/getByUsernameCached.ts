import { getByUsername } from "@/shared/actions/users";
import { cache } from "react";

/**
 * This prevents us from contacting the server on multiple requests in same layout
 *
 * Probably skill and/or api design issues:
 * I just didn't want the user ID to be in the url,
 * it's better and more readable to use his username,
 * but not better as a db primary key, so..
 */
export default cache(getByUsername);

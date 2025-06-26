import { createInertiaApp } from '@inertiajs/react'
import { createRoot } from 'react-dom/client'
import React from 'react'

createInertiaApp({
    resolve: name => {
        //@ts-ignore
        const pages = import.meta.glob('./pages/**/*.tsx', { eager: true })
        return pages[`./pages/${name}.tsx`]
    },
    setup({ el, App, props }) {
        createRoot(el).render(<App {...props} />)
    },
})
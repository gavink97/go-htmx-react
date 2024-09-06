import { NavBar } from '@/assets/navbar/navbar';
import { createRoot } from 'react-dom/client';

const headerRoot = document.getElementById('react-nav');
if (!headerRoot) {
	throw new Error('Could not find element with id react-nav');
}
const headerReactRoot = createRoot(headerRoot);
headerReactRoot.render(NavBar());
